package awsplus

import (
	"fmt"
	"log"

	"github.com/crowdmob/goamz/kinesis"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceKinesisStream() *schema.Resource {
	return &schema.Resource{
		Create: resourceKinesisStreamCreate,
		Read:   resourceKinesisStreamRead,
		Update: nil,
		Delete: resourceKinesisStreamDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"shard_count": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceKinesisStreamCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*kinesis.Kinesis)

	streamName := d.Get("name").(string)
	shardCount := d.Get("shard_count").(int)

	log.Printf("[DEBUG] Creating stream: %#v, shards: %#v", streamName, shardCount)

	if err := client.CreateStream(streamName, shardCount); err != nil {
		return fmt.Errorf("Error creating stream: %s", err)
	}

	d.SetId(d.Get("name").(string))

	return nil
}

func resourceKinesisStreamRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*kinesis.Kinesis)

	description, err := client.DescribeStream(d.Get("name").(string))
	if err != nil {
		if kinesisErr, ok := err.(*kinesis.Error); ok && kinesisErr.Code == "ResourceNotFoundException" {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving stream: %s", err)
	}

	d.Set("name", description.StreamName)
	d.SetId(description.StreamName)

	return nil
}

func resourceKinesisStreamDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*kinesis.Kinesis)

	log.Printf("[INFO] Deleting Stream: %v", d.Id())

	if err := client.DeleteStream(d.Get("name").(string)); err != nil {
		return fmt.Errorf("Error deleting stream: %s", err)
	}

	return nil
}
