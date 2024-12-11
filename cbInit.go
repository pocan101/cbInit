package main

import (
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"errors"
	"github.com/couchbase/gocb/v2"
	"gopkg.in/yaml.v3"
)

type ConnectionDetails struct {
	User      string          `yaml:"user"`
	Password  string          `yaml:"password"`
	URL       string          `yaml:"url"`
	CaCert    CaCertificate   `yaml:"ca_certificate"`
}

type CaCertificate struct {
	Enabled bool   `yaml:"enabled"`
	Name    string `yaml:"name"`
	Content string `yaml:"content"`
}

type BucketSettings struct {
	Name           string `yaml:"name"`
	RamQuotaMb     uint64 `yaml:"ram_quota_mb"`
	BucketType     string `yaml:"bucket_type"`
	NumReplicas    uint32 `yaml:"num_replicas"`
	FlushEnabled   bool   `yaml:"flush_enabled"`
	StorageBackend string `yaml:"storage_backend"`
}

type PreStatements struct {
	QueryName string `yaml:"query_name"`
	N1ql      string `yaml:"n1ql"`
}

type PostStatements struct {
	QueryName string `yaml:"query_name"`
	N1ql      string `yaml:"n1ql"`
}

type Configuration struct {
	ConnectionDetails ConnectionDetails `yaml:"connection_details"`
	Buckets           []BucketSettings  `yaml:"buckets"`
	PreStatements     []PreStatements   `yaml:"pre_ddl_statements"`
	PostStatements    []PostStatements  `yaml:"post_ddl_statements"`
}

// Reads YAML configuration into the Configuration struct
func readConfig(filePath string) (*Configuration, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Configuration
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	// Save CA certificate to file if enabled
	if config.ConnectionDetails.CaCert.Enabled {
		err := os.WriteFile(config.ConnectionDetails.CaCert.Name, []byte(config.ConnectionDetails.CaCert.Content), 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to write CA certificate: %w", err)
		}
	}

	return &config, nil
}

// Create buckets based on configuration


// createOrUpdateBucket handles creating or updating a bucket.
func createOrUpdateBucket(cluster *gocb.Cluster, bucketConfig BucketSettings) error {
    manager := cluster.Buckets()

    settings := gocb.BucketSettings{
        Name:           bucketConfig.Name,
        RAMQuotaMB:     bucketConfig.RamQuotaMb,
        NumReplicas:    bucketConfig.NumReplicas,
        BucketType:     gocb.BucketType(bucketConfig.BucketType),
        FlushEnabled:   bucketConfig.FlushEnabled,
        StorageBackend: gocb.StorageBackend(bucketConfig.StorageBackend),
    }

    // Check if the bucket already exists
    existingBucket, err := manager.GetBucket(bucketConfig.Name, nil)
    if err != nil {
        if errors.Is(err, gocb.ErrBucketNotFound) {
            // Create the bucket if it does not exist
            err := manager.CreateBucket(gocb.CreateBucketSettings{
                BucketSettings:           settings,
                ConflictResolutionType: gocb.ConflictResolutionTypeSequenceNumber,
            }, nil)
            if err != nil {
                return fmt.Errorf("failed to create bucket %s: %w", bucketConfig.Name, err)
            }
            log.Printf("Bucket '%s' created successfully!", bucketConfig.Name)
            return nil
        }
        return fmt.Errorf("failed to check bucket existence: %w", err)
    }

    // Compare storage backend
    if string(existingBucket.StorageBackend) != bucketConfig.StorageBackend {
        return fmt.Errorf(
            "cannot update bucket %s because storage backend differs: existing=%s, new=%s",
            bucketConfig.Name,
            existingBucket.StorageBackend,
            bucketConfig.StorageBackend,
        )
    }

    // Update bucket settings
    err = manager.UpdateBucket(settings, nil)
    if err != nil {
        return fmt.Errorf("failed to update bucket %s: %w", bucketConfig.Name, err)
    }

    log.Printf("Bucket '%s' updated successfully!", bucketConfig.Name)
    return nil
}


// Execute N1QL queries
func executeQueries(cluster *gocb.Cluster, statements []PreStatements) error {
	for _, statement := range statements {
		_, err := cluster.Query(statement.N1ql, nil)
		if err != nil {
			log.Printf("Error executing query '%s': %v", statement.QueryName, err)
			return err
		}
		log.Printf("Successfully executed query '%s'", statement.QueryName)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./cbInit <config.yml>")
		os.Exit(1)
	}

	configPath := os.Args[1]

	// Read configuration
	config, err := readConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	// Initialize Couchbase cluster
	opts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: config.ConnectionDetails.User,
			Password: config.ConnectionDetails.Password,
		},
	}

	if config.ConnectionDetails.CaCert.Enabled {
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM([]byte(config.ConnectionDetails.CaCert.Content)) {
			log.Fatalf("Failed to parse CA certificate")
		}

		opts.SecurityConfig = gocb.SecurityConfig{
			TLSRootCAs: caCertPool,
		}
	}

	cluster, err := gocb.Connect(config.ConnectionDetails.URL, opts)
	if err != nil {
		log.Fatalf("Failed to connect to Couchbase cluster: %v", err)
	}
	defer cluster.Close(nil)

	// Execute pre-DDL statements
	err = executeQueries(cluster, config.PreStatements)
	if err != nil {
		log.Fatalf("Failed to execute pre-DDL queries: %v", err)
	}

	// Create buckets
	for _, bucket := range config.Buckets {
		err := createOrUpdateBucket(cluster, bucket)
		if err != nil {
			log.Fatalf("Failed to create bucket: %v", err)
		}
	}

	// Execute post-DDL statements
	postStatements := make([]PreStatements, len(config.PostStatements))
	for i, ps := range config.PostStatements {
		postStatements[i] = PreStatements{QueryName: ps.QueryName, N1ql: ps.N1ql}
	}

	err = executeQueries(cluster, postStatements)
	if err != nil {
		log.Fatalf("Failed to execute post-DDL queries: %v", err)
	}

	log.Println("cbInit execution completed successfully!")
}

