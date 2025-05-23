# Show-Times E-Commerce Backend Rest API

ShowTimes E-Commerce Backend Rest API is a feature-rich backend solution for E-commerce applications developed using Golang with the Gin web framework. This API is designed to efficiently handle routing and HTTP requests while following best practices in code architecture and dependency management.

## Key Features

- **Clean Code Architecture**: The project follows clean code architecture principles, making it maintainable and scalable.
- **Dependency Injection**: Utilizes the Dependency Injection design pattern for flexible component integration.
- **Compile-Time Dependency Injection**: Dependencies are managed using Wire for compile-time injection.
- **Database**: Leverages PostgreSQL for efficient and relational data storage.
- **AWS Integration**: Integrates with AWS S3 for cloud-based storage solutions.
- **E-Commerce Features**: Implements a wide range of features commonly found in e-commerce applications, including cart management,  order management, wallet, offers, and coupon management.

## Deployment

The application is hosted on AWS EC2 and is served using Nginx, ensuring reliability and scalability.

## API Documentation

For interactive API documentation, Swagger is implemented. You can explore and test the API endpoints in real-time.

## Security

Security is a top priority for the project:

- **OTP Verification**: Twilio API is integrated for OTP verification.
- **Payment Integration**: Razorpay API is used for payment processing.

## Getting Started

To run the project locally, you can follow these steps:

1. Clone the repository.
2. Set up your environment with the required dependencies, including Golang, PostgreSQL,  and Wire.
3. Configure your environment variables (e.g., database credentials, AWS keys, Twilio credentials).
4. Build and run the project.

## Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Wire for Dependency Injection](https://github.com/google/wire)
- [PostgreSQL](https://www.postgresql.org/)
- [AWS S3](https://aws.amazon.com/s3/)
- [Swagger API Documentation](https://swagger.io/)
- [Twilio](https://www.twilio.com/)
- [Razorpay](https://razorpay.com/)

# Environment Variables

Before running the project, you need to set the following environment variables with your corresponding values:

## PostgreSQL

- `DB_HOST`: Database host
- `DB_NAME`: Database name
- `DB_USER`: Database user
- `DB_PORT`: Database port
- `DB_PASSWORD`: Database password

## Admin Security

- `AdminAccessKey`: Acceskey for admin
- `AdminRefreshKey`: Refreshkey for admin

## User Security

- `UserAccessKey`: Acceskey for User
- `UserRefreshKey`: Refreshkey for User

## Twilio

- `DB_AUTHTOKEN`: Twilio authentication token
- `DB_ACCOUNTSID`: Twilio account SID
- `DB_SERVICESID`: Twilio services ID

## AWS

- `AWSRegion`: AWS region
- `AWSAccesskeyID`: AWS access key ID
- `AWSSecretaccesskey`: AWS secret access key

## Razor Pay

- `RazorPay_key_id`: Razor pay key ID
- `RazorPay_key_secret`: Razor pay secret key  

Make sure to provide the appropriate values for these environment variables to configure the project correctly.
