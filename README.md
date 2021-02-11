# `terraform-provider-ksql`
[![CircleCI](https://circleci.com/gh/Mongey/terraform-provider-ksql.svg?style=svg&circle-token=320e9b975067221dd59cc169e83b8faf53ea5062)](https://circleci.com/gh/Mongey/terraform-provider-ksql)

A [Terraform][1] plugin for managing [Confluent KSQL Server][2].

## Contents

- [`terraform-provider-ksql`](#terraform-provider-ksql)
  - [Contents](#contents)
  - [Installation](#installation)
    - [Developing](#developing)
    - [Distribuition](#distribuition)
  - [Provider Configuration](#provider-configuration)
    - [Example](#example)
  - [Resources](#resources)
    - [`ksql_stream`](#ksql_stream)
    - [`ksql_table`](#ksql_table)
    - [`ksql_source_connector`](#ksql_source_connector)
    - [`ksql_sink_connector`](#ksql_sink_connector)

## Installation

Download and extract the [latest
release](/latest) to
your [terraform plugin directory][third-party-plugins] (typically `~/.terraform.d/plugins/`)

### Developing

0. [Install go][install-go]
0. Clone repository
0. Build the provider `make build`
0. Run the tests `make test`

### Distribuition

0. Adjust the [VERSION](VERSION)
0. Run the build
   ```
   make -f Makefile all
   ```

## Provider Configuration

### Example

```hcl
provider "ksql" {
  url = "http://localhost:8083"
}
```

## Resources

### `ksql_stream`

A resource for managing KSQL streams
```hcl
resource "ksql_stream" "actions" {
  name = "vip_actions"
  query = "SELECT userid, page, action
              FROM clickstream c
              LEFT JOIN users u ON c.userid = u.user_id
              WHERE u.level =
              'Platinum';"
}
```

the same with just ksql query string:

```hcl
resource "ksql_stream" "actions" {
  ksql = <<EOF
create stream vip_actions SELECT userid, page, action
              FROM clickstream c
              LEFT JOIN users u ON c.userid = u.user_id
              WHERE u.level =
              'Platinum';
EOF
}
```

### `ksql_table`

A resource for managing KSQL tables
```hcl
resource "ksql_table" "users" {
  name = "users-thing"
  query = "SELECT error_code,
            count(*),
            FROM monitoring_stream
            WINDOW TUMBLING (SIZE 1 MINUTE)
            WHERE  type = 'ERROR'
            GROUP BY error_code;"
  }
}
```

the same with just ksql query string:

```hcl
resource "ksql_table" "users" {
  ksql = <<EOF
create table users-thing SELECT error_code,
            count(*),
            FROM monitoring_stream
            WINDOW TUMBLING (SIZE 1 MINUTE)
            WHERE  type = 'ERROR'
            GROUP BY error_code;
EOF
  }
}
```

### `ksql_source_connector`

```hcl
resource "ksql_source_connector" "jdbcconnector" {
  ksql = <<EOF
CREATE SOURCE CONNECTOR `jdbc-connector` WITH(
    "connector.class"='io.confluent.connect.jdbc.JdbcSourceConnector',
    "connection.url"='jdbc:postgresql://localhost:5432/my.db',
    "mode"='bulk',
    "topic.prefix"='jdbc-',
    "table.whitelist"='users',
    "key"='username');
EOF
  }
}
```

### `ksql_sink_connector`

```hcl
resource "ksql_sink_connector" "docs" {
  ksql = <<EOF
CREATE SINK CONNECTOR DOCS WITH (
   'connector.class'          = 'io.confluent.connect.jdbc.JdbcSinkConnector',
   'connection.url'           = 'jdbc:postgresql://localhost:5432/my_db',
   'connection.user'          = 'pguser',
   'connection.password'      = 'pgpass',
   'tasks.max'                = '1',
   'topics'                   = 'docs_avro',
   'batch.size'               = '3',
   'table.name.format'        = 'docs',
   'auto.create'              = 'true',
   'auto.evolve'              = 'true',
   'delete.enabled'           = 'false',
   'pk.mode'                  = 'record_key',
   'pk.fields'                = 'docId',
   'insert.mode'              = 'UPSERT'
);
EOF
  }
}
```

[install-go]: https://golang.org/doc/install#install
[1]: https://www.terraform.io
[2]: https://www.confluent.io/product/ksql/
