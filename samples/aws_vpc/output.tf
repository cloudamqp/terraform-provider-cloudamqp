// VPC info output
output "cloudamqp_vpc_id" {
  value = data.cloudamqp_vpc_info.vpc_info.id
}

output "cloudamqp_vpc_owner_id" {
  value = data.cloudamqp_vpc_info.vpc_info.owner_id
}

// AWS subnet output
output "aws_vpc_id" {
  value = data.aws_subnet.subnet.vpc_id
}

// AWS peering connection output
output "aws_peering_id" {
  aws_vpc_peering_connection.aws_vpc_peering.id
}

// AWS route table
output "aws_route_table_id" {
  value = data.aws_route_table.route_table.route_table_id
}
