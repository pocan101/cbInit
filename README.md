## **cbInit - Couchbase Bucket Initialization Tool**

### **Overview**

`cbInit` is a tool designed to simplify Couchbase cluster management. It automates the creation and configuration of buckets, execution of N1QL queries, and setup of scopes and collections. The tool is configurable via a YAML file, making it suitable for DevOps workflows and reproducible environments.

------

### **Features**

- Create and update Couchbase buckets with detailed settings.
- Execute pre-defined N1QL queries (DDL or other operations).
- Support for TLS using CA certificates.
- Fully customizable through a YAML configuration file.

------

### **Requirements**

1. Couchbase Server
   - Version 7.x or later.
2. Go Environment
   - Go 1.20+ installed.
3. Dependencies
   - Couchbase Go SDK (`gocb`).
   - YAML library for Go (`gopkg.in/yaml.v3`).

------

### **Installation**

#### **Step 1: Clone the Repository**

```bash

git clone https://github.com/your-repository/cbInit.git
cd cbInit
```

#### **Step 2: Install Dependencies**

Initialize the Go module and fetch dependencies:

```bash

go mod tidy
```

#### **Step 3: Build the Binary**

Compile the Go code:

```bash

go build -o cbInit cbInit.go
```

------

### **Usage**

#### **Command**

```bash

./cbInit <config.yml>
```

#### **Example**

```bash

./cbInit bucket_template.yml
```

```bash
2024/12/11 01:33:29 Successfully executed query 'drop collection  `analytics_bucket`.events.collection1''
2024/12/11 01:33:29 Successfully executed query 'drop scope travel-sample.events'
2024/12/11 01:33:29 Successfully executed query 'drop scope example_bucket.events'
2024/12/11 01:33:29 Successfully executed query 'drop scope analytics_bucket.events'
2024/12/11 01:33:29 Successfully executed query 'insert new airline'
2024/12/11 01:33:30 Bucket 'example_bucket' updated successfully!
2024/12/11 01:33:31 Bucket 'analytics_bucket' updated successfully!
2024/12/11 01:33:34 Successfully executed query 'create index idx_airport_country'
2024/12/11 01:33:44 Successfully executed query 'create index idx_routes_sourceairport'
2024/12/11 01:33:48 Successfully executed query 'create  index idx_hotels_city'
2024/12/11 01:33:48 Successfully executed query 'create  scope `travel-sample`.events'
2024/12/11 01:33:48 Successfully executed query 'create scope `example_bucket`.events'
2024/12/11 01:33:49 Successfully executed query 'create scope `analytics_bucket`.events'
2024/12/11 01:33:49 Successfully executed query 'create collection `travel-sample`.events.collection1''
2024/12/11 01:33:49 Successfully executed query 'create collection `example_bucket`.events.collection1''
2024/12/11 01:33:49 Successfully executed query 'create collection`analytics_bucket`.events.collection1''
2024/12/11 01:33:49 cbInit execution completed successfully!
Process 68948 has exited with status 0
Detaching
```


------

### **Configuration File (YAML)**

The configuration file specifies the connection details, bucket settings, and N1QL queries to execute. Below is an example:

#### **Structure**

```yaml

connection_details:
  user: "Administrator"
  password: "password"
  url: "couchbases://<your-couchbase-host>"
  ca_certificate:
    enabled: true
    name: "ca_certificate.pem"
    content: |
      -----BEGIN CERTIFICATE-----
      ...certificate content...
      -----END CERTIFICATE-----

buckets:
  - name: "example_bucket"
    ram_quota_mb: 5024
    bucket_type: "membase"
    num_replicas: 1
    flush_enabled: false
    storage_backend: "magma"

pre_ddl_statements:
  - query_name: "drop_scope"
    n1ql: "DROP SCOPE `travel-sample`.events"

post_ddl_statements:
  - query_name: "create_index"
    n1ql: "CREATE INDEX idx_country ON `travel-sample`(country) WHERE type = 'airport';"
```

------

### **Features of the Configuration**

1. **Connection Details**:
   - Define credentials, Couchbase URL, and CA certificate settings.
   - CA certificates are written to a local file if provided.
2. **Buckets**:
   - Configure bucket settings, such as name, memory quota, replicas, and storage backend.
   - Automatically creates or updates buckets unless the storage backend differs.
3. **Queries**:
   - **Pre-Statements**: N1QL queries executed before bucket creation.
   - **Post-Statements**: N1QL queries executed after bucket creation.

------

### **Logging**

- The tool provides debug-level logging for detailed troubleshooting.
- Modify the logging configuration in the `enableLogging()` function in the code if needed.

------

### **Error Handling**

1. **Connection Issues**:
   - Ensure the Couchbase server is reachable and TLS settings are correct.
2. **Configuration Errors**:
   - Verify the YAML file structure matches the expected schema.
3. **Bucket Creation**:
   - If a bucket with a different storage backend exists, the tool will fail with an error.

------



# Configuration Template Documentation

The configuration template (`config.yml`) is used to define the Couchbase cluster connection, bucket settings, and N1QL queries to be executed. Below is a detailed explanation of each field.

---

## Top-Level Fields

### `connection_details`
Defines the connection details for the Couchbase cluster.

- **`user`** *(string, required)*  
  The username for authenticating with the Couchbase cluster.

- **`password`** *(string, required)*  
  The password for authenticating with the Couchbase cluster.

- **`url`** *(string, required)*  
  The connection string for the Couchbase cluster, including the protocol (`couchbases://` for TLS or `couchbase://` for non-TLS).

- **`ca_certificate`** *(object, optional)*  
  Provides details for using a custom Certificate Authority (CA) for TLS connections.  
  - **`enabled`** *(boolean, required)*  
    Whether the CA certificate should be used for authentication.  
    - `true`: Enables CA certificate authentication.  
    - `false`: Disables CA certificate authentication.  
  - **`name`** *(string, required if `enabled` is true)*  
    The name of the file where the CA certificate will be written.  
  - **`content`** *(string, required if `enabled` is true)*  
    The PEM-encoded content of the CA certificate.

---

### `buckets`
An array of objects defining the buckets to be created or updated.

Each object includes:

- **`name`** *(string, required)*  
  The name of the bucket.

- **`ram_quota_mb`** *(integer, required)*  
  The RAM quota for the bucket in megabytes.

- **`bucket_type`** *(string, required)*  
  The type of bucket.  
  - `membase`: A Couchbase bucket.  
  - `memcached`: A Memcached bucket.

- **`num_replicas`** *(integer, required)*  
  The number of replicas for the bucket (0 to 3).

- **`flush_enabled`** *(boolean, required)*  
  Whether the bucket allows flushing of data.  
  - `true`: Flushing is allowed.  
  - `false`: Flushing is not allowed.

- **`storage_backend`** *(string, required)*  
  The storage engine used for the bucket.  
  - `couchstore`: Default storage engine.  
  - `magma`: High-performance storage engine for larger datasets.

---

### `pre_ddl_statements`
An array of objects defining pre-DDL (Data Definition Language) statements to be executed before bucket creation.

Each object includes:

- **`query_name`** *(string, required)*  
  A descriptive name for the query.

- **`n1ql`** *(string, required)*  
  The N1QL query to be executed.

---

### `post_ddl_statements`
An array of objects defining post-DDL (Data Definition Language) statements to be executed after bucket creation.

Each object includes:

- **`query_name`** *(string, required)*  
  A descriptive name for the query.

- **`n1ql`** *(string, required)*  
  The N1QL query to be executed.

---

## Example Configuration File

```yaml

version: 1.0
title: Couchbase Template for Bucket Initialization
author: Jose Cortes Diaz
contact: jose.cortesdiaz@couchbase.com
organization: Couchbase

description: >
  This template demonstrates an example use case for creating and configuring 
  two Couchbase buckets, along with initializing example indexes. It serves as 
  a guide for using templates to streamline Couchbase bucket setup.
notes: >
  Ensure that the Couchbase server is running and accessible before applying 
  this template. Customize the buckets and indexes as per your use case.


created: 2024-12-07
last_updated: 2024-12-11

tags:
  - couchbase
  - bucket initialization
  - indexes
  - templates

use_cases:
  - Simplified Couchbase bucket setup
  - Reproducible environment configuration
  - Template-based management for DevOps workflows

compatibility:
  couchbase_version: "7.x and above"
  sdk_version: "3.x and above"
requirements:
  - A Couchbase server instance
  - Administrative access to Couchbase
  - Couchbase C++ SDK installed


usage:
    - command: cbInit --template bucket_template.yml
    - description: Use this command to apply the bucket initialization template.
documentation: README.txt

connection_details:
  user: Administrator
  password: password
  url: couchbases://ec2-3-88-72-54.compute-1.amazonaws.com
  ca_certificate:
    enabled: true
    name: "ca_certificate.pem" # Name of the generated file
    content: |
      -----BEGIN CERTIFICATE-----
      MIIDDDCCAfSgAwIBAgIIGA4YNVtbAwYwDQYJKoZIhvcNAQELBQAwJDEiMCAGA1UE
      AxMZQ291Y2hiYXNlIFNlcnZlciAzNzJlZWY3MTAeFw0xMzAxMDEwMDAwMDBaFw00
      OTEyMzEyMzU5NTlaMCQxIjAgBgNVBAMTGUNvdWNoYmFzZSBTZXJ2ZXIgMzcyZWVm
      NzEwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDeMWK3CWJcbxXqD818
      8JMIEz689ByxQpClCTqQn3IRLpQufV5EWOKoeQrIrLoxfBHy93j0LOaR57eztdXr
      8R7SN4TNwK/hocjpbtjTic2nzAecMAQqhs8mc9aPNDtbJvxS/0IyxPnyEK0dN1ib
      gUKX38NBW5uiAU5bugU+4mZX/KMFuI+DMvnKqItANB+Q7Opcwtna1Ke121zFntWV
      z0TYReFgW4lgirvoC0gxmi61E0Jtr0ZXSMLv0L2SP5kLelKqxjYU4tmgdJvrrHD1
      1pnVwhUrA6XpXy+PYehEUrRHPIa+j/3cv7ng/2D/XY7DS3kMxe4GzWQ89cKatcYZ
      q+3TAgMBAAGjQjBAMA4GA1UdDwEB/wQEAwIBhjAPBgNVHRMBAf8EBTADAQH/MB0G
      A1UdDgQWBBR8JCP7Kll4RS8Eh6awFFrdS1sGFjANBgkqhkiG9w0BAQsFAAOCAQEA
      yK0TKRTxZXMHFAJ3JBCHOipVPlbCtaRBjCh6UKchcNXZESdA2BbDjlJWwP19x8xG
      sBvpp1GWTkMzuRNilsWBterCnUW/HjPRagAXUZ6dx9gCZoRacMzEfKid6rz+1Drj
      kX2syG3b0rfHcQqc/CMzrLruk+mLobHTxB5CRjtIJoIOMn7+laLwwnC+luZm5Z66
      wGJMaY7BagzZk4GrWQXoeq2IdJyPEX2gnMOljpS7QqyodBseXw0u6RrohX0dhZX2
      nzFSHWXJHYCYY9iYAFjJ8dUpffPBxYMwDHsz7XJwPooa+ylTXXcAGQb/irpEAb3S
      uiXs9TY6QIDez+pxkDH0fw==
      -----END CERTIFICATE-----


buckets:
  - name: example_bucket
    ram_quota_mb: 5024
    bucket_type: membase # couchbase or memcached
    num_replicas: 1
    flush_enabled: false
    storage_backend: magma # couchstore or magma

  - name: analytics_bucket
    ram_quota_mb: 200
    bucket_type: membase
    num_replicas: 1
    flush_enabled: true
    storage_backend: couchstore

pre_ddl_statements:

 - query_name: drop index idx_airport_country
   n1ql: 'DROP INDEX `travel-sample`.idx_airport_country'
 - query_name: drop index idx_routes_sourceairport
   n1ql: 'DROP INDEX `travel-sample`.idx_routes_sourceairport'
 - query_name: drop index idx_hotels_city
   n1ql: 'DROP INDEX `travel-sample`.idx_hotels_city'
 - query_name: drop collection `travel-sample`.events.collection1'
   n1ql: 'DROP COLLECTION `travel-sample`.events.collection1'
 - query_name: drop collection `example_bucket`.events.collection1'
   n1ql: 'DROP COLLECTION `example_bucket`.events.collection1'
 - query_name: drop collection  `analytics_bucket`.events.collection1'
   n1ql: 'DROP COLLECTION `analytics_bucket`.events.collection1'
 - query_name: drop scope travel-sample.events
   n1ql: 'DROP SCOPE `travel-sample`.events'

 - query_name: drop scope example_bucket.events
   n1ql: 'DROP SCOPE `example_bucket`.events'

 - query_name: drop scope analytics_bucket.events
   n1ql: 'DROP SCOPE `analytics_bucket`.events'
 - query_name: insert new airline
   n1ql: 'UPSERT INTO `travel-sample` (KEY, VALUE) VALUES ("hotel_9999", {"type": "hotel", "name": "Dreamland Hotel", "address": "123 Paradise Lane", "city": "San Francisco", "state": "CA", "country": "USA", "phone": "123-456-7890", "url": "http://dreamlandhotel.com"});'
 
post_ddl_statements:

  - query_name: create index idx_airport_country
    n1ql: 'CREATE INDEX idx_airport_country ON `travel-sample`(country, airportname, city) WHERE type = "airport";'

  - query_name: create index idx_routes_sourceairport
    n1ql: 'CREATE INDEX idx_routes_sourceairport ON `travel-sample`(sourceairport, destinationairport, airlineid) WHERE type = "route";'

  - query_name: create  index idx_hotels_city
    n1ql: 'CREATE INDEX idx_hotels_city ON `travel-sample`(city, name, address, country) WHERE type = "hotel";'

  - query_name: create  scope `travel-sample`.events
    n1ql: 'CREATE SCOPE `travel-sample`.events'

  - query_name: create scope `example_bucket`.events
    n1ql: 'CREATE SCOPE `example_bucket`.events'

  - query_name: create scope `analytics_bucket`.events
    n1ql: 'CREATE SCOPE `analytics_bucket`.events'
  - query_name: create collection `travel-sample`.events.collection1'
    n1ql: 'CREATE COLLECTION `travel-sample`.events.collection1'

  - query_name: create collection `example_bucket`.events.collection1'
    n1ql: 'CREATE COLLECTION `example_bucket`.events.collection1'

  - query_name: create collection`analytics_bucket`.events.collection1'
    n1ql: 'CREATE COLLECTION `analytics_bucket`.events.collection1'
```



### **Compatibility**

- Couchbase Server: Version 7.x and above.
- Couchbase SDK: Version 3.x and above.

------



### **Contact**

For questions or support, reach out to: **Jose Cortes Diaz**
Email: jose.cortesdiaz@couchbase.com