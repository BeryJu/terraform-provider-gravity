# Records are identified by zone, hostname, type and uid, separated by colons
terraform import gravity_dns_record.example 'my-domain.com.:www:A:0'
