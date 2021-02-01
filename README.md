# pvdf

## Overview

pvdf stand for 'PersistentVolume Disk Free'. The idea is to provide a quick an useful view of all disk usage on a Kubernetes cluster hosting some PersistentVolume.

Here is a sample output:

```
NAMESPACE	NODE	PV NAME			        POD NAME		        REQ.	STORAGE CLASS	SIZE	FREE	%USED
gha-1			    datalake1					                    50Gi			        ???	???	???
gha-1			    datalake1-pv		    gha2posix-161210..zq	50Gi			        49Gi	25Gi	48%
kluster1	s1	    pvc-46885e83-6d1..81	kluster1-kafka-2	    20Gi	topolvm-ssd	    19Gi	0	    100%
kluster1	s1	    pvc-6d582c36-628..67	kluster1-zookeeper-2	2Gi	    topolvm-ssd	    2036Mi	1996Mi	1%
kluster1	s2	    pvc-9a9c9050-aa0..be	kluster1-kafka-1	    20Gi	topolvm-ssd	    19Gi	0	    100%
kluster1	s2	    pvc-f6e7a949-50d..f7	kluster1-zookeeper-1	2Gi	    topolvm-ssd	    2036Mi	1996Mi	1%
kluster1	s3	    pvc-1dad433b-901..fe	kluster1-kafka-0	    20Gi	topolvm-ssd	    19Gi	0	    100%
kluster1	s3	    pvc-ba88bd67-eea..20	kluster1-zookeeper-0	2Gi	    topolvm-ssd	    2036Mi	1996Mi	1%
minio1		s1	    pvc-3142c088-818..1b	minio1-zone-0-1		    5Gi	    topolvm-ssd	    5105Mi	5072Mi	0%
minio1		s1	    pvc-94478c98-0ba..9c	minio1-zone-0-1		    5Gi	    topolvm-ssd	    5105Mi	5072Mi	0%
minio1		s2	    pvc-1843ac08-cf4..84	minio1-zone-0-2		    5Gi	    topolvm-ssd	    5105Mi	5072Mi	0%
minio1		s2	    pvc-a1a73eac-1cd..08	minio1-zone-0-2		    5Gi	    topolvm-ssd	    5105Mi	5072Mi	0%
minio1		s3	    pvc-22c38e92-646..02	minio1-zone-0-0		    5Gi	t   opolvm-ssd	    5105Mi	5072Mi	0%
minio1		s3	    pvc-3d42d920-eaa..48	minio1-zone-0-0		    5Gi	    topolvm-ssd	    5105Mi	5072Mi	0%
minio1		s3	    pvc-b157dcf9-3d5..e4	minio1-zone-0-3		    5Gi	    topolvm-ssd	    5105Mi	5072Mi	0%
minio1		s3	    pvc-bf77d019-42a..3f	minio1-zone-0-3		    5Gi	    topolvm-ssd	    5105Mi	5072Mi	0%
prometheus	s1	    pvc-c97bf74c-b7b..c6	prometheus-prome..-0	10Gi    topolvm-hdd	    10220Mi	0	    100%
```

One can see on this sample this cluster is unhealthy, with some disk full.

## Usage


## Architecture


# Installation 


