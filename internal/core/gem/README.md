# Core: GEM

The **G**lobal **E**ntity **M**anager is what controls the storage, lifetime,
sorting and all other aspects of your games [entities](../entities/README.md).

It can either be used directly for simple applications or with an intermediary
layer like [scenes](../../modules/scenes/README.md) that helps with grouping
and managing multiple entities at a time.

In the following, only direct usage will be shown, visit the
[scenes](../../modules/scenes/README.md) section for information about how to
use scenes.

## How does it work?

The real entity storage in the GEM is a 2D matrix, or a list of lists.
Each top level list represents a layer group, and contains a list of entities
within that layer group.
This way of storing entities lends itself to quickly retrieving the layers for
rendering, but isn't an intuitive interface to work with.

For this reason, the GEM also stores some hierarchical data parallel to the
layer-entity matrix. This builds a tree structure, where each entity has a list
of children associated with it. The branches in this tree share certain
properties, such as their transform (position, rotation, scale). They also
share a lifetime, which is a thing the GEM controls as well.

Every entity must expose functions like `Init()`, `Deinit()`, `OnChildAdded(c)`, ...
The GEM automatically calls these functions whenever an entity is added or
removed from its control. When the parent of an entity is removed, all children
are automatically removed as well.

## Usage

First, we need an [entity](../entities/README.md) to manage:
```go
myEntity := entities.NewTestEntity()
```

We then add the entity to them GEM. we can do this before or during our game
loop.
```go
gem.AddEntity(gem.GetRoot(), myEntity, gem.DefaultLayer)
```

We must then call the `Update()` and `Draw()` functions on that entity every
frame. We can do that either manually per entity, or let the GEM do it:

```go
for !shouldStop {
    // ...

    // manually
    myEntity.Update()

    // or automatically for all entities in that layer.
    gem.Update(gem.GetByLayer(gem.DefaultLayer))

    // notice that we could also have called Update/Draw on ALL entities, not just that layer:
    gem.Draw(gem.GetAll()) // automatically for all entities in all layers.

    // ...
}
```

> TIP: Filtering and Drawing by Layer works well in conjunction with render
> stages from the [render](../render/README.md) module, should you need
> multiple render layers.

To remove entities, we simply call:

```go
gem.Remove(myEntity)
```
