variable "vcpus" {}
variable "memory" {}

provider "libvirt" {
  uri = "qemu:///system"
}

resource "libvirt_volume" "<VolumeName>" {
  name   = "<VolumeName>"
  pool   = "<PoolName>"
  source = "<Source>"
  format = "<Format>"
}

resource "libvirt_domain" "<DomainName>" {
  name   = "<DomainName>"
  memory = var.memory
  vcpu   = var.vcpus

  network_interface {
      network_name = "<NetworkName>"
      macvtap      = true
      addresses    = ["<IPAddress>"]
  }

  disk {
      volume_id = libvirt_volume.<VolumeName>.id
  }
}

provider "wireguard" {}

resource "wireguard_interface" "<InterfaceName>" {
  private_key       = wireguard_private_key.<KeyName>.private_key_base64
}

resource "wireguard_peer" "<PeerName>" {
}
