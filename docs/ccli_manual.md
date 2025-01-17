## Catalog Command Line Interface

We present an overview and user manual for the Catalog Command Line interface (ccli) utility. The ccli is used to perform certain data access and change 
operations on the Software Parts Catalog (SPC) such as adding new and updating parts and part profiles, uploading archives, and retrieving part and profile data.
The following operations are supported:

- **add** part <file.yml> - adds a new part record to the catalog. See the 'add' section below for the format the required yml file. For example:
```
$ ccli add part openssl-1.1.1n.yml
```
- **add** profile <file.yml> - adds a new part profile document to the catalog. See the 'add' section below for the format the required yml file. For example:
```
$ ccli add profile profile_openssl-1.1.1n.yml
```
- **query** <string> - enables one to query the catalog for part data. For example:
```
$ ccli query '...'
```
- **export** 
part fvc <file_verification_code>| sha256 \<Sha256>| id <catalog_id> -o <file.yml>
Export out data for a given part. 
```
  $ ccli export part id sdl3ga-naTs42g5-rbow2A -o file.yml
```
- **export** 
template <part | security | quality | licensing> -o <Path.yaml>
Export template for part or profile
```
ccli export template security -o file.yml

ccli export template license -o file.yml
```
- **update** <file.yml> - enables one to update selective data fields of a part record. See the 'update' section below for the format the 
required yml file. For example:
```
$ ccli update openssl-1.1.1n.v4.yml
```
-  **upload** <source archive> - uploads the specified source archive. A a new part record will be created if it does not correspond part record exists otherwise
it will be associated with an existing part if it already exists.  
```
$ ccli upload openssl-1.1.1n.tar.gz
```
- **find** 
part \<query> - searches catalog for matching part names and displays corresponding data
id \<catalog_id> - retrieves a part from catalog using id
sha256 \<sha256> - returns part id using given sha256
fvc \<file_verification_code> - returns part id using given file verification code
```
$ ccli find part busybox
$ ccli find sha256 <sha256>
```
- **find**
profile <security|quality|licensing> <catalog_id> - retrieves a profile from the catalog based on type and part id.
```
ccli find profile security werS12-da54FaSff-9U2aef
```
- **delete**
 <catalog_id> - deletes a part from the catalog using part id if the part has no related parts. Recursive flag can be used to delete a part and its sub-parts as long as they have no other related parts.
```
ccli delete adjb23-A4D3faTa-d95Xufs
```
```
ccli delete adjb23-A4D3faTa-d95Xufs --recursive
```

## Add
- ### Part
```
format: 1.0
fvc: "4656433200c41848db861f590cd5cb929265011204d6ea4851f966fd5f4a33295a2569b35f"
sha256: "faeeb244c35a348a334f4a59e44626ee870fb07b6884d68c10ae8bc19f83a694"
catalog_id: null
name: "busybox"
version: "1.35.0"
type: "archive"
family_name: "busybox"
label: "busybox-1.35.0"
description: |
  Provides many of the Unix utilities in a single executable file.
license: 
  license_expression: "GPL-2.0"
  analysis_type: "automation/topline/1.1"
size: null
aliases: 
  - "busybox-1.35.0"
  - "busybox-1.35.0.r3"
comprised_of: null
composite_list: null
```
- ### Security Profile
```
profile: 'security'
format: 1.0
name: 'busybox'
version: "1.35.0"
fvc: "4656433200c41848db861f590cd5cb929265011204d6ea4851f966fd5f4a33295a2569b35f"
sha256: 'faeeb244c35a348a334f4a59e44626ee870fb07b6884d68c10ae8bc19f83a694' 
catalog_id: null
## ----------------------------------------------------
##                 CVEs
## ----------------------------------------------------
cve_list:
  ## cve
  - cve_id: 'CVE-2022-30065'
    description: |
      A use-after-free in Busybox 1.35-x's awk applet leads to denial of service and 
      possibly code execution when processing a crafted awk pattern in the copyvar function.
    status: 'Open'
    date: '2022-05-02'
    link: 'https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-30065'
    comments: |
      patch is available
    references: []
  ## -----------------
  ## cve
  ## -----------------
  - cve_id: 'CVE-2022-28391'
    description: |
      BusyBox through 1.35.0 allows remote attackers to execute arbitrary code if netstat 
      is used to print a DNS PTR record's value to a VT compatible terminal. Alternatively, 
      the attacker could choose to change the terminal's colors.
    status: 'open'
    date: '2022-04-03'
    comments: null
    link: 'https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-28391'
    references: []
```
- ### Licensing Profile
```
profile: "licensing"
format: 1.0
name: "busybox"
version: "1.35.0"
fvc: "4656433200c41848db861f590cd5cb929265011204d6ea4851f966fd5f4a33295a2569b35f"
sha256: "faeeb244c35a348a334f4a59e44626ee870fb07b6884d68c10ae8bc19f83a694"
catalog_id: null

## ----------------------
## License Analysis
## ----------------------
license_analysis:
  - license_expression: GPL-2.0
    analysis_type: expert/mark.gisi@windriver.com
    comments: null
  - license_expression: GPL-2.0
    analysis_type: automation/topline/1.0
    comments: null

## ----------------------
## Copyrights
## ----------------------
copyrights:
  - Copyright 1999-2005 Erik Andersen

## ----------------------
## Legal Notices
## ----------------------
legal_notice: |
  Copyright 1999-2005 Erik Andersen

  This package is free software; you can redistribute it and/or modify
  it under the terms of the GNU General Public License as published by
  the Free Software Foundation; version 2 dated June, 1991.

  This package is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with this package; if not, write to the Free Software
  Foundation, Inc., 51 Franklin St, Fifth Floor, Boston,
  MA 02110-1301, USA.

## ----------------------
## Other Legal Notices
## ----------------------
other_legal_notices: null
```
- ### Quality Profile
```
profile: "quality"
format: 1.0
name: busybox
version: "1.35.0"
fvc: "4656433200c41848db861f590cd5cb929265011204d6ea4851f966fd5f4a33295a2569b35f"
sha256: "faeeb244c35a348a334f4a59e44626ee870fb07b6884d68c10ae8bc19f83a694"
catalog_id: null
## ----------------------------------------------------
##                 Bugs
## ----------------------------------------------------
bug_list:
  ## -----------------
  ## bug
  ## -----------------
  - name: "Bug 14376"
    id: "14376"
    description: |
      Tar component has a memory leak bug when trying to unpack a tar file..
    status: "Open"
    level: "P5"
    date: "2022-11-23"
    link: https://bugs.busybox.net/show_bug.cgi?id=14376
    comments: null
    references:
      - "[Busybox Bug Report](https://bugs.busybox.net/show_bug.cgi?id=14376)"
  ## -----------------
  ## bug
  ## -----------------
  - name: "Bug 14536"
    id: "14536"
    description: |
      Awk from busybox-v1.35.0 doesn’t work.
    status: "Open"
    level: "P5"
    date: "2022-01-21"
    link: https://bugs.busybox.net/show_bug.cgi?id=14536
    comments: |
      There is a patch that reverts awk to busybox-1.33.1.
    references:
      - "[Busybox Bug Report](https://bugs.busybox.net/show_bug.cgi?id=14536)"
```

## CCLI Example Usage
```
$ ccli examples
    $ ccli add part openssl-1.1.1n.yml
    $ ccli add profile profile_openssl-1.1.1n.yml
    $ ccli query "{part(id:\"aR25sd-V8dDvs2-p3Gfae\"){file_verification_code}}"
    $ ccli export part id sdl3ga-naTs42g5-rbow2A -o file.yml
    $ ccli export template security -o file.yml
    $ ccli update openssl-1.1.1n.v4.yml
    $ ccli upload openssl-1.1.1n.tar.gz
    $ ccli find part busybox
    $ ccli find sha256 2493347f59c03...
    $ ccli find profile security werS12-da54FaSff-9U2aef
    $ ccli delete adjb23-A4D3faTa-d95Xufs
    $ ccli ping
```

## Updating License Example
### Steps:
**(1) Export out main profile:**
```
$ ccli export part sha256 12cec6bd2b16d8a9446dd16130f2b92982f1819f6e1c5f5887b6db03f5660d28 -o busybox-33.yml
Part successfully exported to path: busybox-33.yml
```
Two ways to obtain the sha256:
  - use the linux coomand: sha256sum busybox-33.1.tar.gz
  - look up the part in the catalog and copy the sha256 from one of the archives  listed.
    
**(2) Edit file: busybox-33.yml and change the license to GPL-2.0:**
```
license:
    license_expression: ""
    analysis_type: ""
```
To:
```
license:
    license_expression: "GPL-2.0"
    analysis_type: "human-analysis/mark.gisi@windriver.com"
```
**(3) Perform an update:**
```
ccli $ ccli update busybox-33.yml

Part successfully updated
{
  "ID": "858f8261-19dc-441a-9cfb-5fe2e38345ac",
  "PartType": "/archive",
  "ContentType": "tbd",
  "Version": "1.33.1",
  "Name": "busybox",
  "Label": "busybox-1.33.1",
  "FamilyName": "",
  "FileVerificationCode": "465643320019459ec42edd3ba8723fe40f261c4489c7ed497888777b8396f406dcaa588157",
  "Size": 9925270,
  "License": "GPL-2.0",
  "LicenseRationale": "human-analysis/mark.gisi@windriver.com",
  "Description": "",
  "HomePage": "",
  "Comprised": "00000000-0000-0000-0000-000000000000",
  "Aliases": []
}
```

