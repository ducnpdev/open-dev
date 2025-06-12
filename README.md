# Open-Dev

A comprehensive collection of Go (Golang) development examples, patterns, and best practices. This repository serves as a learning resource and reference implementation for various Go development scenarios.

## Table of Contents

- [Overview](#overview)
- [Getting Started](#getting-started)
- [Core Components](#core-components)
  - [Web Development](#web-development)
  - [Message Queues](#message-queues)
  - [Database & Caching](#database--caching)
  - [Cloud & Infrastructure](#cloud--infrastructure)
  - [Security & Cryptography](#security--cryptography)
  - [Performance & Optimization](#performance--optimization)
  - [Design Patterns & Best Practices](#design-patterns--best-practices)
- [Contributing](#contributing)
- [Contact & Support](#contact--support)
- [License](#license)

## Overview

This repository contains practical examples and implementations of various Go development concepts, from basic web applications to advanced distributed systems patterns. It's designed to help developers learn and implement Go best practices in real-world scenarios.

## Getting Started

1. Clone the repository:
```bash
git clone https://github.com/ducnpdev/open-dev.git
cd open-dev
```

2. Install dependencies:
```bash
go mod download
```

3. Explore the examples in each directory based on your interests.

## Core Components

### Web Development
- **Let's Go**: Build web applications in Go
  - [Basic web application examples](letgo/)
- **Gin Web Framework**: 
  - [RESTful API examples](gin-web-framework/)
  - [Middleware implementations](gin-web-framework/)
  - [Error handling patterns](gin-web-framework/)

### Message Queues
- **Kafka**: 
  - [Advanced Consumer Patterns](kafka/README_VI.md#consumer-advanced)
    - Consumer group management
    - Offset handling
    - Error recovery
  - [Advanced Producer Patterns](kafka/README_VI.md#topic-advanced)
    - Topic management
    - Partition strategies
    - Message delivery guarantees
  - [Basic Examples](kafka/)
    - [Installation guide](kafka/)
    - Simple producer/consumer setup
    - Configuration examples
- **RabbitMQ**: 
  - [Code examples and patterns](https://github.com/ducnpdev/golang-rabbitmq)
  - Message routing
  - Queue management
  - Exchange types

### Database & Caching
- **Redis**:
  - [Rate Limiting Implementation](redis/)
    - Token bucket algorithm
    - Sliding window implementation
    - [Rate limit article](https://viblo.asia/p/golang-ratelimit-la-gi-su-dung-ByEZkn3qKQ0)
  - [Pub/Sub Patterns](redis/pubsub/)
    - Publisher implementation
    - Subscriber patterns
    - Message handling
  - [Installation Guide](redis/)
    - Setup instructions
    - Configuration examples
    - Connection management
- **Database Connection Management**: 
  - [Strategy patterns for database connections](https://viblo.asia/p/golang-chien-luoc-quan-ly-ket-noi-database-thong-qua-gorm-zXRJ8r6dVGq)
  - Connection pooling
  - Retry mechanisms
  - Error handling

### Cloud & Infrastructure
- **AWS**:
  - [Lambda Functions](aws/lambda/):
    - [Simple function examples](aws/lambda/)
    - [CRUD operations with Postgres](aws/lambda/crud)
    - [Docker integration guide](https://aws.amazon.com/blogs/compute/migrating-aws-lambda-functions-from-the-go1-x-runtime-to-the-custom-runtime-on-amazon-linux-2/)
  - [S3 Operations](aws/s3/):
    - File upload/download
    - Bucket management
    - Access control
  - [Textract Integration](aws/textract/):
    - [Document processing](aws/textract/main.go)
    - [Identity verification](https://docs.aws.amazon.com/textract/index.html)
- **Terraform**:
  - [SQS Setup](terraform/apps/sqs/):
    - Queue creation
    - Access policies
    - Monitoring setup
  - [API Gateway Configuration](terraform/apps/apigw/):
    - Endpoint setup
    - Integration patterns
    - Security configurations
- **Serverless Framework**: 
  - [Implementation guide](https://viblo.asia/p/golang-cai-dat-serverless-framework-lambada-aws-4dbZNQXkKYM)
  - Function deployment
  - Environment management

### Security & Cryptography
- **RSA**: 
  - [Encryption/Decryption examples](usecase/rsa/REAME.md)
  - Key generation
  - Message signing
- **ECDSA**: 
  - [Digital signature implementation](usecase/ecdsa/ecdsa.go)
  - Key pair management
  - Signature verification
- **TLS**: 
  - [Certificate handling](tls/base64parse/)
  - Base64 conversion
  - Key management
- **Password Security**: 
  - Hashing algorithms
  - Salt implementation
  - Secure storage practices

### Performance & Optimization
- **HTTP Optimization**: 
  - [Connection reuse patterns](http/reuse-http/)
  ```go
  // Example of HTTP connection reuse
  func Reuse() {
      clientTrace := &httptrace.ClientTrace{
          GotConn: func(info httptrace.GotConnInfo) {
              log.Printf("conn was reused: %t", info.Reused)
          },
      }
      // ... implementation details in http/reuse-http/
  }
  ```
- **Performance Testing**: 
  - [Standard function benchmarks](performances/)
  - [Error handling performance](https://opendev.hashnode.dev/golang-test-performance-function-standard-1)
  - [String conversion benchmarks](https://opendev.hashnode.dev/golang-test-performance-function-standard-1)
- **Concurrency Patterns**: 
  - [Essential Go concurrency patterns](concurrency/patterns/)
  - Channel management
  - Goroutine best practices

### Design Patterns & Best Practices
- **UUID Usage**: 
  - Benefits:
    - Distributed system compatibility
    - Unpredictability for security
  - Drawbacks:
    - Storage overhead
    - Non-sequential nature
- **HTTPS Implementation**: 
  - TLS handshake process
  - Certificate management
  - Security best practices
- **SSO (Single Sign-On)**: 
  - Authentication flow
  - Token management
  - Session handling
- **Database Design**: 
  - Password storage best practices
  - Salt implementation
  - Hash algorithms

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## Contact & Support

- **Author**: Đức Nguyễn Phúc
- **Email**: ducnp09081998@gmail.com or phucducktpm@gmail.com
- **Social Media**:
  - [Facebook](https://www.facebook.com/phucducdev/)
  - [LinkedIn](https://www.linkedin.com/in/phucducktpm/)
  - [Hashnode](https://hashnode.com/@OpenDev)
  - [Telegram](https://t.me/OpenDevGolang)

## License

This project is licensed under the terms of the license included in the repository. See the [LICENSE](LICENSE) file for details.

---

⭐ If you find this repository helpful, please give it a star! Your support helps maintain and improve this resource for the Go community.