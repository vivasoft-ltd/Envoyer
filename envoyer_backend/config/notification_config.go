package config

import (
	"time"
)

/*
PublisherConfig -
When you want a single message to be delivered to a single queue, you can
publish to the default exchange with the routingKey of the queue name.  This is
because every declared queue gets an implicit route to the default exchange.

Since publishings are asynchronous, any undeliverable message will get returned
by the server.  Add a listener with Channel.NotifyReturn to handle any
undeliverable message when calling publish with either the mandatory or
immediate parameters as true.

Publishings can be undeliverable when the mandatory flag is true and no queue is
bound that matches the routing key, or when the immediate flag is true and no
consumer on the matched queue is ready to accept the delivery.

This can return an error when the channel, connection or socket is closed.  The
error or lack of an error does not indicate whether the server has received this
publishing.

It is possible for publishing to not reach the broker if the underlying socket
is shut down without pending publishing packets being flushed from the kernel
buffers.  The easy way of making it probable that all publishings reach the
server is to always call Connection.Close before terminating your publishing
application.  The way to ensure that all publishings reach the server is to add
a listener to Channel.NotifyPublish and put the channel in confirm mode with
Channel.Confirm.  Publishing delivery tags and their corresponding
confirmations start at 1.  Exit when all publishings are confirmed.

When Publish does not return an error and the channel is in confirm mode, the
internal counter for DeliveryTags with the first confirmation starts at 1.
*/
type PublisherConfig struct {
	ServerUrl      string
	ExchangeConfig ExchangeConfig
	Mandatory      bool
	Immidiate      bool
	Exclusive      bool
	ContentType    string
	BackoffPolicy  []time.Duration
}

/*
ConsumerConfig -
Consume immediately starts delivering queued messages.

Begin receiving on the returned chan Delivery before any other operation on the
Connection or Channel.

Continues deliveries to the returned chan Delivery until Channel.Cancel,
Connection.Close, Channel.Close, or an AMQP exception occurs.  Consumers must
range over the chan to ensure all deliveries are received.  Unreceived
deliveries will block all methods on the same connection.

All deliveries in AMQP must be acknowledged.  It is expected of the consumer to
call Delivery.Ack after it has successfully processed the delivery.  If the
consumer is cancelled or the channel or connection is closed any unacknowledged
deliveries will be requeued at the end of the same queue.

The consumer is identified by a string that is unique and scoped for all
consumers on this channel.  If you wish to eventually cancel the consumer, use
the same non-empty identifier in Channel.Cancel.  An empty string will cause
the library to generate a unique identity.  The consumer identity will be
included in every Delivery in the ConsumerTag field

When autoAck (also known as noAck) is true, the server will acknowledge
deliveries to this consumer prior to writing the delivery to the network.  When
autoAck is true, the consumer should not call Delivery.Ack. Automatically
acknowledging deliveries means that some deliveries may get lost if the
consumer is unable to process them after the server delivers them.
See http://www.rabbitmq.com/confirms.html for more details.

When exclusive is true, the server will ensure that this is the sole consumer
from this queue. When exclusive is false, the server will fairly distribute
deliveries across multiple consumers.

The noLocal flag is not supported by RabbitMQ.

It's advisable to use separate connections for
Channel.Publish and Channel.Consume so not to have TCP pushback on publishing
affect the ability to consume messages, so this parameter is here mostly for
completeness.

When noWait is true, do not wait for the server to confirm the request and
immediately begin deliveries.  If it is not possible to consume, a channel
exception will be raised and the channel will be closed.

Optional arguments can be provided that have specific semantics for the queue
or server.

Inflight messages, limited by Channel.Qos will be buffered until received from
the returned chan.

When the Channel or Connection is closed, all buffered and inflight messages will
be dropped.

When the consumer tag is cancelled, all inflight messages will be delivered until
the returned chan is closed.
*/
type ConsumerConfig struct {
	ExchangeConfig ExchangeConfig
	QueueConfig    QueueConfig
	ServerUrl      string
	NoWait         bool
	AutoAck        bool
	BackoffPolicy  []time.Duration
}

/*
ExchangeConfig declares an exchange on the server. If the exchange does not
already exist, the server will create it.  If the exchange exists, the server
verifies that it is of the provided type, durability and auto-delete flags.

Each exchange belongs to one of a set of exchange kinds/types implemented by
the server. The exchange types define the functionality of the exchange - i.e.
how messages are routed through it. Once an exchange is declared, its type
cannot be changed.  The common types are "direct", "fanout", "topic" and
"headers".

Durable and Non-Auto-Deleted exchanges will survive server restarts and remain
declared when there are no remaining bindings.  This is the best lifetime for
long-lived exchange configurations like stable routes and default exchanges.

Non-Durable and Auto-Deleted exchanges will be deleted when there are no
remaining bindings and not restored on server restart.  This lifetime is
useful for temporary topologies that should not pollute the virtual host on
failure or after the consumers have completed.

Non-Durable and Non-Auto-deleted exchanges will remain as long as the server is
running including when there are no remaining bindings.  This is useful for
temporary topologies that may have long delays between bindings.

Durable and Auto-Deleted exchanges will survive server restarts and will be
removed before and after server restarts when there are no remaining bindings.
These exchanges are useful for robust temporary topologies or when you require
binding durable queues to auto-deleted exchanges.

Note: RabbitMQ declares the default exchange types like 'amq.fanout' as
durable, so queues that bind to these pre-declared exchanges must also be
durable.

Exchanges declared as `internal` do not accept publishings. Internal
exchanges are useful when you wish to implement inter-exchange topologies
that should not be exposed to users of the broker.

When noWait is true, declare without waiting for a confirmation from the server.
The channel may be closed as a result of an error.  Add a NotifyClose listener
to respond to any exceptions.
*/
type ExchangeConfig struct {
	Name       string
	Type       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
}

/*
QueueConfig declares a queue to hold messages and deliver to consumers.
Declaring creates a queue if it doesn't already exist, or ensures that an
existing queue matches the same parameters.

Every queue declared gets a default binding to the empty exchange "" which has
the type "direct" with the routing key matching the queue's name.  With this
default binding, it is possible to publish messages that route directly to
this queue by publishing to "" with the routing key of the queue name.

The queue name may be empty, in which case the server will generate a unique name
which will be returned in the Name field of Queue struct.

Durable and Non-Auto-Deleted queues will survive server restarts and remain
when there are no remaining consumers or bindings.  Persistent publishings will
be restored in this queue on server restart.  These queues are only able to be
bound to durable exchanges.

Non-Durable and Auto-Deleted queues will not be redeclared on server restart
and will be deleted by the server after a short time when the last consumer is
canceled or the last consumer's channel is closed.  Queues with this lifetime
can also be deleted normally with QueueDelete.  These durable queues can only
be bound to non-durable exchanges.

Non-Durable and Non-Auto-Deleted queues will remain declared as long as the
server is running regardless of how many consumers.  This lifetime is useful
for temporary topologies that may have long delays between consumer activity.
These queues can only be bound to non-durable exchanges.

Durable and Auto-Deleted queues will be restored on server restart, but without
active consumers will not survive and be removed.  This Lifetime is unlikely
to be useful.

Exclusive queues are only accessible by the connection that declares them and
will be deleted when the connection closes.  Channels on other connections
will receive an error when attempting  to declare, bind, consume, purge or
delete a queue with the same name.

When noWait is true, the queue will assume to be declared on the server.  A
channel exception will arrive if the conditions are met for existing queues
or attempting to modify an existing queue from a different connection.

When the error return value is not nil, you can assume the queue could not be
declared with these parameters, and the channel will be closed.
*/
type QueueConfig struct {
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
}
