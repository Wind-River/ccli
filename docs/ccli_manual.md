## Catalog Command Line Interface

We present an overview and user manual for the Catalog Command Line interface (ccli) utiltity. The ccli is used to perform certain data access and change 
operations on the Software Parts Catalog (SPC) such as adding new and updating parts and part profiles, uploading archives, and retrieving part and profile data.
The following operations are supported:

- **add** --part <file.yml> - adds a new part record to the catalog. See the 'add' section below for the format the required yml file. For example:
```
$ ccli add --part openssl-1.1.1n.yml
```
- **add** --profile <file.yml> - adds a new part profile document to the catalog. See the 'add' section below for the format the required yml file. For example:
```
$ ccli add --profile profile_openssl-1.1.1n.yml
```
- **query** <string> - enables one to query the catalog for part data. For example:
```
$ ccli query '...'
```
- **export** --part label|<fvc>|sha256|catlog_id|  -o <file.yml>. Export out data for a given part. 
```
  $ ccli export --part label|<fvc>|sha256|catlog_id|  -o <file.yml>
```
- **update** --part <file.yml> - enables one to update selective data fields of a part record. See the 'update' section below for the format the 
required yml file. For example:
```
$ ccli update --part openssl-1.1.1n.v4.yml
```
- **update** --profile <file.yml> - enables one to update selective data fields of a part profile document. See the 'update' section below for the
format of the required yml file. For example:
```
$ ccli update --profile profile_openssl-1.1.1n.v2.yml
```
-  **upload** <source archive> - uploads the specified source archive. A a new part record will be created if it does not correspond part record exists otherwise
it will be associated with an existing part if it already exists.  
```
$ ccli upload openssl-1.1.1n.tar.gz
```
-  **upload** <source archive> - uploads the specified source archive. A a new part record will be created if it does not correspond part record exists otherwise
it will be associated with an existing part if it already exists.  
```
$ ccli upload openssl-1.1.1n.tar.gz
```
