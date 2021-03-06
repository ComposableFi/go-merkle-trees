# Go Merkle Trees

This is the most advanced Merkle tree library for Go. Basic features
include building a Merkle tree, creation, and verification of Merkle proofs for
single and several elements, i.e. multi-proofs. Advanced features include making
transactional changes to the tree and rolling back to any previously committed
tree state, similarly to Git.

The library is highly customizable. Hashing function and the way how the tree
is built can be easily configured through a special trait.

## About Merkle trees

Merkle trees, also known as hash trees, are used to verify that two or more
parties have the same data without exchanging the entire data collection.

Merkle trees are used in Git, Mercurial, ZFS, IPFS, Bitcoin, Tendermint, Ethereum, Cassandra,
and many more. In Git, for example, Merkle trees are used to find a delta
between the local and remote repository states to transfer only the difference
between them over the network. In Bitcoin, Merkle trees are used to verify that
a transaction was included in the block without downloading the whole block
contents. ZFS uses Merkle trees to quickly verify data integrity, offering
protection from silent data corruption caused by phantom writes, bugs in disk
firmware, power surges, and other causes.

## About Merkle Mountain Range

A type of merkle tree that can be visualized as many (perfect) merkle trees which are then combined into 1,
by creating a single root from all of their peaks. The rules for making the tree(s) 
however are rigidly deterministic such that the entire structure depends only on 
the number of items put into it. When appending leaves, nodes are not updated, only appended.
This makes for a minimal amount of total hash computations (~2n), 
while maintaining the property that no nodes are updated when appending new leaves (nodes are merely added). 
Because of this, the merkle inclusion proofs are also only appended to meaning that later proofs are sufficient (supersets of) earlier proofs.

## Usage

Get the library:

```
go get -u github.com/ComposableFi/go-merkle-trees
```

simply import merkle package and follow the examples
```
import "github.com/ComposableFi/go-merkle-trees/merkle"
```

and for mmr
```
import "github.com/ComposableFi/go-merkle-trees/mmr"
```

### Hasher
You need to implement the hasher type of your desired hashing mechanism, the hasher type should implement the Hash method with this signature:
```
Hash(data []byte) ([]byte, error)
```


## Examples

[Sha256](https://github.com/ComposableFi/go-merkle-trees/tree/main/examples/sha256)

[Keccak](https://github.com/ComposableFi/go-merkle-trees/tree/main/examples/keccak)


## Reference
- Inspired from [rs-merkle](https://github.com/antouhou/rs-merkle) algoritms and tests.
- Converted test cases from [snowbridge](https://github.com/Snowfork/snowbridge/blob/main/ethereum/test/test_merkle_tree.js).