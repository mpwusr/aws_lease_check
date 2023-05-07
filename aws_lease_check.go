package AWS_lease_check

import (
	"fmt"
	"time"
)

package main

import (
"fmt"
"time"

"github.com/aws/aws-sdk-go/aws"
"github.com/aws/aws-sdk-go/aws/session"
"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	// Initialize AWS session
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	// Specify the name of the DynamoDB table
	tableName := "my-tokens-table"

	// Specify the TTL attribute name
	ttlAttributeName := "expiration"

	// Specify the lease duration in seconds
	leaseDuration := 60

	// Start an infinite loop to periodically check the leases
	for {
		// Calculate the current time plus the lease duration
		expirationTime := time.Now().Add(time.Duration(leaseDuration) * time.Second)

		// Create a DynamoDB query input
		input := &dynamodb.QueryInput{
			TableName: aws.String(tableName),
			KeyConditions: map[string]*dynamodb.Condition{
				ttlAttributeName: {
					ComparisonOperator: aws.String("LT"),
					AttributeValueList: []*dynamodb.AttributeValue{
						{
							N: aws.String(fmt.Sprintf("%d", expirationTime.Unix())),
						},
					},
				},
			},
		}

		// Query the DynamoDB table for expired tokens
		result, err := svc.Query(input)
		if err != nil {
			fmt.Println("Error querying DynamoDB:", err)
			continue
		}

		// Process the expired tokens
		for _, item := range result.Items {
			// Do something with the expired token
			fmt.Println("Expired token:", item["token"].N)
		}

		// Sleep for the lease duration before checking again
		time.Sleep(time.Duration(leaseDuration) * time.Second)
	}
}
