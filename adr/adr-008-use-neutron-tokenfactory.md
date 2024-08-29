# ADR 008: Use Neutron Token Factory by copying code into Mantrachain

## Context

The token factory is a notoriously forked and fragmented module.  To use it we needed to choose between:

* Osmosis (which used sdk 47)
* Neutron (which uses sdk 50)
* Eve (which was missing a patch in Neutron's)

And between weather or not we will:

* import the code from another repository
* add the code to mantrachain by attributed copy paste

## Decision

We will use Neutron's token factory and memorialize the commit we fork from.  This should be the smoothest path forward for us because it won't need to be upgraded to 50, or patched.

Importing Neutron left undesirable changes to our go.mod file, so we will copy the code in and attribute Neutron.

## Status

Accepted
