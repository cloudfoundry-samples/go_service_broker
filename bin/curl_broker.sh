curl -X GET http://username:password@localhost:8001/v2/catalog

curl -X PUT http://username:password@localhost:8001/v2/service_instances/instance_id-111 -d '{
  "service_id":"service-guid-111",
  "plan_id":"plan-guid-111",
  "organization_guid": "org-guid-111",
  "space_guid":"space-guid-111",
  "parameters": {"ami_id":"ami-ecb68a84"}
}' -H "Content-Type: application/json"

curl -X GET http://username:password@localhost:8001/v2/service_instances/instance_id-111

curl -X PUT http://username:password@localhost:8001/v2/service_instances/instance_id-111/service_bindings/binding_id-111 -d '{
  "plan_id":        "plan-guid-here",
  "service_id":     "service-guid-here",
  "app_guid":       "app-guid-here"
}' -H "Content-Type: application/json"

curl -X DELETE http://username:password@localhost:8001/v2/service_instances/instance_id-111/service_bindings/binding_id-111

curl -X DELETE http://username:password@localhost:8001/v2/service_instances/instance_id-111