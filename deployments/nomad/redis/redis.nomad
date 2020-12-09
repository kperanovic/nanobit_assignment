job "redis" {
  region = "global"
  datacenters = ["dc1"]
  type = "service"

  update {
    max_parallel = 1
    min_healthy_time = "10s"
    healthy_deadline = "3m"
    
    auto_revert = false
    canary = 0
  }

  group "cache" {
    count = 1

    restart {
      attempts = 10
      interval = "5m"
      delay = "10s"
      mode = "delay"
    }

    task "redis" {
      driver = "docker"

      config {
        image = "redis:3.2"
        port_map {
          db = 6379
        }
      }
      
      resources {
        cpu    = 500
        memory = 256
        network {
          mbits = 10
          port "db" {}
        }
      }

      # Setting service config for service discovery system (Consul)
      service {
        name = "global-redis-check"
        tags = ["global", "cache", "urlprefix-/redis" ]
        port = "db"
        check {
          name     = "alive"
          type     = "tcp"
          interval = "10s"
          timeout  = "2s"
        }
      }
    }
  }
}