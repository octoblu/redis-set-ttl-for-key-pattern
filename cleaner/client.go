package cleaner

import "github.com/garyburd/redigo/redis"

// Client connects to redis an sets ttls for
// all keys matching an expression
type Client struct {
	conn    redis.Conn
	pattern string
}

// New constructs a new cleaner.Client
func New(pattern, redisURI string) (*Client, error) {
	conn, err := redis.DialURL(redisURI)
	if err != nil {
		return nil, err
	}

	return &Client{conn, pattern}, nil
}

// Clean will scan redis and set tlls on
// all records returned. It returns the cursor
// for the next call
func (client *Client) Clean(cursor int) (int, error) {
	conn := client.conn

	scanResults, err := redis.Values(conn.Do("SCAN", cursor, "MATCH", client.pattern, "COUNT", 1000))
	if err != nil {
		return -1, err
	}

	newCursor, _ := redis.Int(scanResults[0], nil)
	keys, _ := redis.Strings(scanResults[1], nil)
	for _, key := range keys {
		_, err = conn.Do("EXPIRE", key, 86400)
		if err != nil {
			return -1, err
		}
	}

	return newCursor, nil
}

// Close closes the redis connection
func (client *Client) Close() error {
	return client.conn.Close()
}
