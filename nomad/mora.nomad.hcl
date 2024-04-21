
job "mora" {

    datacenters = ["dc1"]
    type = "service"

    group "blockchain" {
        count = 1

        network {
            mode = "bridge"

            port "api" {}

            port "rpc" {}

            port "metrics" {
                to = 2112
            }
        }

        service {
            name = "mora-api" 
            port = "api"

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
            name = "mora-metrics"
            port = "metrics"

            tags = [
                "metrics"
            ]

            provider = "consul"
        }

        service {
            name = "mora-rpc"
            port = "rpc"

            provider = "consul"

            connect {
                sidecar_service {}
            }
        }

        volume "kv-volume" {
            type = "host"
            source = "mora-kv-volume"
            read_only = false
        }

        task "server" {
            driver = "docker"

            config {
                image = "trevatk/mora:v0.0.1"
                ports = [ "api", "rpc", "metrics" ]
            }

            volume_mount {
                volume = "kv-volume"
                destination = "/var/lib/mora/kv"
                read_only = false
            }

            env {
                SERVER_HTTP_PORT = "${NOMAD_PORT_api}"
                SERVER_GRPC_PORT = "${NOMAD_PORT_rpc}"
                KV_DIR = "/var/lib/mora/kv"
            }
        }
    }
}