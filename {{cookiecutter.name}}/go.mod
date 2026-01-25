module github.com/{{ cookiecutter.repo_owner }}/{{ cookiecutter.repo_name }}

go 1.23

require (
	github.com/conductorone/baton-sdk v0.7.1
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
	go.uber.org/zap v1.27.0
	google.golang.org/protobuf v1.36.3
)
