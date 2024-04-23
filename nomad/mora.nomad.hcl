
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

        template {
            data = << EOH
            server {
                bind_addr = "{{ env 'NOMAD_UPSTREAM_ADDR_rpc' }}"
                advertise_addr = ""

                ports {
                    http = "${NOMAD_PORT_api}"
                    grpc = "${NOMAD_PORT_rpc}"
                }
            }

            raft {}
            EOH
            destination = "/etc/mora/config.hcl"
        }

        task "server" {
            driver = "docker"

            config {
                image = "trevatk/mora:v0.0.1"
                ports = [ "api", "rpc", "metrics" ]
            }

            env {
                SERVER_HTTP_PORT = "${NOMAD_PORT_api}"
                SERVER_GRPC_PORT = "${NOMAD_PORT_rpc}"
            }
        }
    }
}