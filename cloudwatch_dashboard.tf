# CloudWatch Dashboard for EggCarton Project
resource "aws_cloudwatch_dashboard" "eggcarton_dashboard" {
  dashboard_name = "EggCarton-Monitoring"

  dashboard_body = jsonencode({
    widgets = [
      # DynamoDB ItemCount - SingleValue Widget
      {
        type = "metric"
        properties = {
          metrics = [
            ["AWS/DynamoDB", "ItemCount", "TableName", aws_dynamodb_table.egg_carton.name, { stat = "Average", label = "Total Items" }]
          ]
          view    = "singleValue"
          region  = var.aws_region
          title   = "DynamoDB Item Count"
          period  = 300
          yAxis = {
            left = {
              showUnits = false
            }
          }
        }
        width  = 6
        height = 3
        x      = 0
        y      = 0
      },
      # Lambda Invocations vs Errors - Stacked Area Chart
      {
        type = "metric"
        properties = {
          metrics = [
            ["AWS/Lambda", "Invocations", "FunctionName", aws_lambda_function.put_egg.function_name, { stat = "Sum", label = "Put Egg Invocations" }],
            [".", ".", ".", aws_lambda_function.get_egg.function_name, { stat = "Sum", label = "Get Egg Invocations" }],
            [".", ".", ".", aws_lambda_function.break_egg.function_name, { stat = "Sum", label = "Break Egg Invocations" }],
            [".", "Errors", ".", aws_lambda_function.put_egg.function_name, { stat = "Sum", label = "Put Egg Errors" }],
            [".", ".", ".", aws_lambda_function.get_egg.function_name, { stat = "Sum", label = "Get Egg Errors" }],
            [".", ".", ".", aws_lambda_function.break_egg.function_name, { stat = "Sum", label = "Break Egg Errors" }]
          ]
          view    = "timeSeries"
          stacked = true
          region  = var.aws_region
          title   = "Lambda Invocations vs Errors"
          period  = 300
          yAxis = {
            left = {
              label     = "Count"
              showUnits = false
            }
          }
        }
        width  = 12
        height = 6
        x      = 6
        y      = 0
      },
      # Lambda Duration - Line Chart (Average and p99)
      {
        type = "metric"
        properties = {
          metrics = [
            ["AWS/Lambda", "Duration", "FunctionName", aws_lambda_function.put_egg.function_name, { stat = "Average", label = "Put Egg Avg" }],
            ["...", { stat = "p99", label = "Put Egg p99" }],
            [".", ".", ".", aws_lambda_function.get_egg.function_name, { stat = "Average", label = "Get Egg Avg" }],
            ["...", { stat = "p99", label = "Get Egg p99" }],
            [".", ".", ".", aws_lambda_function.break_egg.function_name, { stat = "Average", label = "Break Egg Avg" }],
            ["...", { stat = "p99", label = "Break Egg p99" }]
          ]
          view    = "timeSeries"
          stacked = false
          region  = var.aws_region
          title   = "Lambda Duration (Average and p99)"
          period  = 300
          yAxis = {
            left = {
              label     = "Milliseconds"
              showUnits = false
            }
          }
        }
        width  = 18
        height = 6
        x      = 0
        y      = 6
      },
      # Log Insights - Security Alerts Table
      {
        type = "log"
        properties = {
          query   = <<-EOT
            SOURCE '/aws/lambda/${aws_lambda_function.put_egg.function_name}'
            | SOURCE '/aws/lambda/${aws_lambda_function.get_egg.function_name}'
            | SOURCE '/aws/lambda/${aws_lambda_function.break_egg.function_name}'
            | fields @timestamp, @message, @logStream
            | filter @message like /Security Alert/
            | sort @timestamp desc
            | limit 10
          EOT
          region  = var.aws_region
          title   = "Recent Security Alerts"
          stacked = false
        }
        width  = 24
        height = 6
        x      = 0
        y      = 12
      }
    ]
  })
}

# Output the dashboard URL
output "cloudwatch_dashboard_url" {
  description = "URL to view the EggCarton CloudWatch Dashboard"
  value       = "https://console.aws.amazon.com/cloudwatch/home?region=${var.aws_region}#dashboards:name=${aws_cloudwatch_dashboard.eggcarton_dashboard.dashboard_name}"
}
