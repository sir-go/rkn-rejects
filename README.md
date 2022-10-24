## DNS/DPI sniffers and NFT-tables rules

The parental control project contains four utilities to get white and black lists from RKN service and 
completely isolate one certain host from denied resources. Utilities installed at the router between
the host and uplink.

### Utilities

 - [get_rkn](cmd/get_rkn)

    SOAP-client for [RKN service](https://vigruzki.rkn.gov.ru/services/OperatorRequest/?wsdl), 
    gets b/w lists, parses them and fills the redis cache.
    
 - [dns-sniffer](cmd/dns-sniffer)

    watches all DNS traffic, collects A-records, and checks resolved hosts by b/w lists.
    
    If the record in the answer is not denied then it puts to the white IP list with resolved TTL.

 - [dpi-sniffer](cmd/dpi-sniffer)

    sniffs all traffic and for denied IP or hostname detection.

 - [check](cmd/check)

    the testing tool for regularly checking the quality of other parts working. It gets the 
    accessibility of hosts by their list in Redis, runs goroutines with HTTP-clients, and collects status codes.

    It can be run from the docker container for routing all traffic through the router's firewall.

### NF-tables

Traffic to sniffers redirects by the `nf_queue` kernel module. All traffic rejects by default 
except DNS requests and answers. NF-tables rules have a list of allowed IP addresses. 
Every record in the list has a TTL and deletes when this time is expired.
Also, there are the lists filled manually and imported to the rules.

Utilities control the firewall with internal library [internal/fw](internal/fw).

Predefined useful nftables configs located in [configs/nftables](configs/nftables)
