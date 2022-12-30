# lambda simple:
- create simple lambda function
- read blog: https://viblo.asia/p/golang-aws-lambda-thong-qua-serverless-framework-phan-1-EoW4ob9xVml
- handle request and response in api: https://viblo.asia/p/golang-aws-lambda-thong-qua-serverless-framework-phan-2-5pPLkPY8VRZ

## Handle API CRUD
- Create table user với script:
```
CREATE TABLE "users" (
    "id" bigserial,
     user_name character varying(50) COLLATE pg_catalog."default",
    name character varying(50) COLLATE pg_catalog."default" NOT NULL,
    phone character varying(50) COLLATE pg_catalog."default",
    PRIMARY KEY ("id")
);
```

- Add new user:
