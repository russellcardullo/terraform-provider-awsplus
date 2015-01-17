package awsplus

import (
	"fmt"
	"log"

	"github.com/crowdmob/goamz/sqs"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSQSQueue() *schema.Resource {
	return &schema.Resource{
		Create: resourceSQSQueueCreate,
		Read:   resourceSQSQueueRead,
		Update: nil,
		Delete: resourceSQSQueueDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"url": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSQSQueueCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AWSClient).sqsconn

	queueName := d.Get("name").(string)

	log.Printf("[DEBUG] Creating queue: %#v", queueName)

	queue, err := client.CreateQueue(queueName)
	if err != nil {
		return fmt.Errorf("Error creating queue: %s", err)
	}

	d.Set("url", queue.Url)
	d.SetId(queue.Url)

	return nil
}

func resourceSQSQueueRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AWSClient).sqsconn

	queue, err := client.GetQueue(d.Get("name").(string))
	if err != nil {
		if queueErr, ok := err.(*sqs.Error); ok && queueErr.Code == "ResourceNotFoundException" {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving queue: %s", err)
	}

	d.Set("url", queue.Url)
	d.SetId(queue.Url)

	return nil
}

func resourceSQSQueueDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AWSClient).sqsconn

	log.Printf("[INFO] Deleting Queue: %v", d.Id())

	queue, err := client.GetQueue(d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("Error retrieving queue: %s", err)
	}

	if _, err := queue.Delete(); err != nil {
		return fmt.Errorf("Error deleting queue: %s", err)
	}

	return nil
}
