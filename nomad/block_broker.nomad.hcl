
job "message-broker" {

    datacenters = [""]
    type = "service"

    namespace = ""

    group "" {
        count = 1

        network {
            mode = "bridge"

            port "dashboard" {}

            port "messenger" {}
        }

        service {
            name = "message-broker" 
            port = "dashboard"

            tags = [
                "traefik.enable=true",
            ]

            provider = "consul"

            connect {
                sidecar_service {}
            }
        }

        service {
            name = "broker-messenger"
            port = "messenger"

            provider = "consul"

            connect {
                sidecar_service {}
            }
        }

        volume "kv-volume" {
            type = "host"
            source = "broker-volume"
            read_only = false
        }

        task "server" {
            driver = "docker"

            config {
                image = ""
                ports = ["dashboard", "messenger"]
            }

            volume {
                volume = "kv-volume"
                destination = "/var/lib/broker/kv"
                read_only = false
            }

            env {
                SERVER_HTTP_PORT = "${NOMAD_PORT_dashboard}"
                SERVER_GRPC_PORT = "${NOMAD_PORT_messenger}"
                KV_DIR = "/var/lib/broker/kv"
            }
        }
    }
}