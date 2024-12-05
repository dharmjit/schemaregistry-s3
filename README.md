# SchemaRegistryS3

## What is Schema Registry
Schema Registry primarily targets the lifecycle of data contracts in various data exchange formats(`JSON`, `AVRO`, `ProtoBuf`) among producers and consumers. This means it provides functionality around storing,validating, retreival and evolution of the data contracts/schemas.

### Popular Schema Registry choices
- Confluent Schema Registry
- Karapace
- Apicurio Registry
- Redpanda Schema Registry

## Why another Schema Registry
Most of the schema registry options are either too tightly integrated with the vendor complete offering ecosystem(for e.g. Confluent Schema Registry with Confluent Platform/Confluent Cloud) or for the persistent storage they again requires kafka, SQL Databases, etc. which makes it a heavy dependency.

So we thought why not give it a try to build a simple library for schema registry backed by S3 Compatible Object Storage. We also discovered that there are less similar OSS options out there.

## What is SchemaRegistryS3
It's a Go library to support storing, validating and retrieving data schemas. The data schemas will be stored in S3 storage in directory-like hierarchy using keys so that's its also easier to navigate them from the service UI. 

```
schemas/
├── eventtype1/
│   ├── v1/
│   │   └── eventtype1-v1.json
│   ├── v2/
│       └── eventtype1-v2.json
├── eventtype2/
│   ├── v1/
│       └── eventtype2-v1.json
```