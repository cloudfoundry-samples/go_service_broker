curl -X GET http://username:password@localhost:8001/v2/catalog

curl -X PUT http://username:password@localhost:8001/v2/service_instances/instance_guid-111 -d '{
  "service_id":"service-guid-111",
  "plan_id":"plan-guid",
  "organization_guid": "org-guid",
  "space_guid":"space-guid",
  "parameters": {"ami_id":"ami-ecb68a84"}
}' -H "Content-Type: application/json"

curl -X PUT http://username:password@localhost:8001/v2/service_instances/instance_guid-222 -d '{
  "service_id":"service-guid-222",
  "plan_id":"plan-guid",
  "organization_guid": "org-guid",
  "space_guid":"space-guid",
  "parameters": {}
}' -H "Content-Type: application/json"

curl -X GET http://username:password@localhost:8001/v2/service_instances/instance_guid-111
curl -X GET http://username:password@localhost:8001/v2/service_instances/instance_guid-222

curl -X PUT http://username:password@localhost:8001/v2/service_instances/instance_guid-111/service_bindings/binding_guid-111 -d '{
  "plan_id":        "plan-guid",
  "service_id":     "service-guid-111",
  "app_guid":       "app-guid"
}' -H "Content-Type: application/json"

curl -X PUT http://username:password@localhost:8001/v2/service_instances/instance_guid-222/service_bindings/binding_guid-222 -d '{
  "plan_id":        "plan-guid",
  "service_id":     "service-guid-222",
  "app_guid":       "app-guid"
}' -H "Content-Type: application/json"

curl -X DELETE http://username:password@localhost:8001/v2/service_instances/instance_guid-111/service_bindings/binding_guid-111
curl -X DELETE http://username:password@localhost:8001/v2/service_instances/instance_guid-111

curl -X DELETE http://username:password@localhost:8001/v2/service_instances/instance_guid-222/service_bindings/binding_guid-222
curl -X DELETE http://username:password@localhost:8001/v2/service_instances/instance_guid-222
