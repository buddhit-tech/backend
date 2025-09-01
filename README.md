# School-Auth


curl http://localhost:8080/healthz
curl -X POST http://localhost:8080/student/login -H "Content-Type: application/json" -d '{"student_id":"S001"}'
curl -X POST http://localhost:8080/student/otp/verify -H "Content-Type: application/json" -d '{"uid":"S001","otp":"123456"}'
curl http://localhost:8080/students
curl http://localhost:8080/teachers
