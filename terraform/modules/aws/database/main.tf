terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }
}

# Environment-specific defaults
locals {
  instance_class_defaults = {
    dev     = "db.t3.micro"
    staging = "db.t3.small"
    prod    = "db.r6g.large"
  }

  storage_defaults = {
    dev     = 20
    staging = 50
    prod    = 100
  }

  instance_class      = var.instance_class != "" ? var.instance_class : local.instance_class_defaults[var.env_type]
  allocated_storage   = var.allocated_storage > 0 ? var.allocated_storage : local.storage_defaults[var.env_type]
  multi_az            = var.env_type == "prod" ? true : false
  backup_retention    = var.env_type == "prod" ? var.backup_retention_period : 7
  deletion_protection = var.env_type == "prod" ? true : false

  db_name = replace("${var.app_name}_${var.env_type}", "-", "_")
}

# Generate random password for database
resource "random_password" "db_password" {
  length  = 32
  special = true
  # Exclude characters that might cause issues in connection strings
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

# Store password in AWS Secrets Manager
resource "aws_secretsmanager_secret" "db_password" {
  name        = "${var.app_name}-${var.env_type}-db-password"
  description = "Database password for ${var.app_name} ${var.env_type} environment"

  tags = merge(
    var.tags,
    {
      Name        = "${var.app_name}-${var.env_type}-db-password"
      App_Name    = var.app_name
      Env_Type    = var.env_type
      ManagedBy   = "devplatform-cli"
      Timestamp   = timestamp()
    }
  )
}

resource "aws_secretsmanager_secret_version" "db_password" {
  secret_id     = aws_secretsmanager_secret.db_password.id
  secret_string = jsonencode({
    username = "dbadmin"
    password = random_password.db_password.result
    engine   = "postgres"
    host     = aws_db_instance.main.address
    port     = aws_db_instance.main.port
    dbname   = local.db_name
  })
}

# Create DB subnet group
resource "aws_db_subnet_group" "main" {
  name       = "${var.app_name}-${var.env_type}-db-subnet-group"
  subnet_ids = var.subnet_ids

  tags = merge(
    var.tags,
    {
      Name      = "${var.app_name}-${var.env_type}-db-subnet-group"
      App_Name  = var.app_name
      Env_Type  = var.env_type
      ManagedBy = "devplatform-cli"
      Timestamp = timestamp()
    }
  )
}

# Create security group for RDS
resource "aws_security_group" "rds" {
  name        = "${var.app_name}-${var.env_type}-rds-sg"
  description = "Security group for RDS instance"
  vpc_id      = var.vpc_id

  ingress {
    description     = "PostgreSQL from application security groups"
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = var.security_group_ids
  }

  egress {
    description = "Allow all outbound traffic"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(
    var.tags,
    {
      Name      = "${var.app_name}-${var.env_type}-rds-sg"
      App_Name  = var.app_name
      Env_Type  = var.env_type
      ManagedBy = "devplatform-cli"
      Timestamp = timestamp()
    }
  )
}

# Create DB parameter group
resource "aws_db_parameter_group" "main" {
  name   = "${var.app_name}-${var.env_type}-pg-params"
  family = "postgres15"

  parameter {
    name  = "log_connections"
    value = "1"
  }

  parameter {
    name  = "log_disconnections"
    value = "1"
  }

  parameter {
    name  = "log_duration"
    value = "1"
  }

  tags = merge(
    var.tags,
    {
      Name      = "${var.app_name}-${var.env_type}-pg-params"
      App_Name  = var.app_name
      Env_Type  = var.env_type
      ManagedBy = "devplatform-cli"
      Timestamp = timestamp()
    }
  )
}

# Create RDS instance
resource "aws_db_instance" "main" {
  identifier     = "${var.app_name}-${var.env_type}-db"
  engine         = "postgres"
  engine_version = var.engine_version
  instance_class = local.instance_class

  allocated_storage     = local.allocated_storage
  max_allocated_storage = local.allocated_storage * 2
  storage_type          = "gp3"
  storage_encrypted     = true

  db_name  = local.db_name
  username = "dbadmin"
  password = random_password.db_password.result

  db_subnet_group_name   = aws_db_subnet_group.main.name
  vpc_security_group_ids = [aws_security_group.rds.id]
  parameter_group_name   = aws_db_parameter_group.main.name

  multi_az               = local.multi_az
  publicly_accessible    = false
  backup_retention_period = local.backup_retention
  backup_window          = "03:00-04:00"
  maintenance_window     = "mon:04:00-mon:05:00"
  deletion_protection    = local.deletion_protection
  skip_final_snapshot    = var.env_type != "prod"
  final_snapshot_identifier = var.env_type == "prod" ? "${var.app_name}-${var.env_type}-final-snapshot-${formatdate("YYYY-MM-DD-hhmm", timestamp())}" : null

  enabled_cloudwatch_logs_exports = ["postgresql", "upgrade"]

  tags = merge(
    var.tags,
    {
      Name      = "${var.app_name}-${var.env_type}-db"
      App_Name  = var.app_name
      Env_Type  = var.env_type
      ManagedBy = "devplatform-cli"
      Timestamp = timestamp()
    }
  )
}
