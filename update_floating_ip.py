import digitalocean
import os

manager = digitalocean.Manager()

domain = manager.get_domain(os.getenv("DOMAIN"))
record = [rec for rec in domain.get_records() if rec.type == 'A'][0]
ip = manager.get_floating_ip(record.data)

if not ip.droplet:
    print("Ip not assigned.")
    k8_workers = manager.get_all_droplets(tag_name="k8s:worker")
    k8_worker = k8_workers[0]
    ip.assign(k8_worker.id)
