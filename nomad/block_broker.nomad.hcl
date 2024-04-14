
job "message-broker" {

    datacenters = ["dc1"]
    type = "service"

    namespace = "structx"

    group "trevatk" {
        count = 1

        network {
            mode = "bridge"

            port "dashboard" {}

            port "rpc" {}

            port "metrics" {
                to = 2112
            }
        }

        service {
            name = "message-broker-dashboard" 
            port = "dashboard"

            tags = [
                "traefik.enable=true",
            ]

            provider = "consul"

            connect {
                sidecar_service {}
            }

            check {
                type = "http"
                path = "/health"
                interval = "10s"
                timeout = "3s"
            }
        }

        service {
            name = "message-broker-metrics"
            port = "metrics"

            tags = [
                "metrics"
            ]

            provider = "consul"
        }

        service {
            name = "broker-messenger-rpc"
            port = "rpc"

            provider = "consul"

            connect {
                sidecar_service {}
            }
        }

        volume "kv-volume" {
            type = "host"
            source = "block-broker-volume"
            read_only = false
        }

        task "server" {
            driver = "docker"

            config {
                image = "trevatk/message-broker:v0.0.1"
                ports = [ "dashboard", "rpc", "metrics" ]
            }

            volume_mount {
                volume = "kv-volume"
                destination = "/var/lib/broker/kv"
                read_only = false
            }

            env {
                SERVER_HTTP_PORT = "${NOMAD_PORT_dashboard}"
                SERVER_GRPC_PORT = "${NOMAD_PORT_rpc}"
                KV_DIR = "/var/lib/broker/kv"
            }
        }
    }
}