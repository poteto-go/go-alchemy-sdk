---
name: unit test expert
description: Expertise in ut.
  Use when write to ut or refactoring ut.
  Except e2e testing `e2e`.
---

# Outline

- Rules
  - Details on rule
- Tips for ut

# Rules

- AMAA pattern
- Easy test case(suite) name for domain expert
- Minimal mock usage
- Increase refactoring resistance
- Minimal test code
- Fast

## AMAA Pattern

Arrange, Mock(if exists), Act, Assert Pattern

EX)

```go
func Test_BuyItemAndReduceInventory(t *testing.T) {
  // Arrange
  store := Store{}
  store.SetItems("Shampoo", 10)

  // Mock
  // ...mocking stripe

  // Act
  inventory := store.Buy("Shampoo", 6, 100)

  // Assert
  assert.Equal(inventory, 4)
}
```

## Easy test case(suite) name for domain expert

### Anti Pattern

- Name the parts to be tested, such as function names and class names.
  EX) `TestEther_GetBlockNumber`

- Mention detailed logic in a white box manner
  EX) `it should call eth_blockNumber and return result`

### Best Practice

- Name based on behavior
  EX) `Test_GetBlockNumberForBlockChain`
  EX) `it returns blocknumber of blockchain node`

## Minimal mock usage

### Anti Pattern

- Write detailed logic mocks
  EX)
  if impl as below,

  ```go
  func (s *Some) DoSomething() {
    return helper.DoSomething()
  }
  ```

  in test

  ```go
  func Test_DoSomething(t *tesiting.T) {
    helper.mock(func(){
      return "hello"
    })

    s := Some{}
    result := s.DoSomething()
    assert.Equal(t, result, "hello")
  }
  ```

### BestPractice

- Testing detailed logic
  EX)

  ```go
  func Test_DoSomething(t *tesiting.T) {
    // Don't mock detail logic

    s := Some{}
    result := s.DoSomething()
    assert.Equal(t, result, "hello") // assert of detail logic
  }
  ```

- Mock external service behavior
  EX)

  ```go
  func Test_getBalanceOfProvidedAccount(t *testing.T) {
    // mock block chain rpc node
    mock := alchemymock.NewAlchemyHttpMock(setting, t)
  defer mock.DeactivateAndReset()
    mock.RegisterResponder("eth_getBalance", `{"jsonrpc":"2.0","id":1,"result":"0x1234"}`)

    balance, err := alchemy.Core.GetBalance("0x", "latest")
  if err != nil {
  	t.Fatal(err)
  }
  assert.Equal(t, "4660", balance.String())
  }
  ```

## Increase refactoring resistance

### Anti Pattern

- assert in mocking function call
  EX)

  ```go
  func Test_1(...) {
    // Mock
    // ... mock `func1`

    // Act
    result := Do()

    // Assert
    assert.Equal(t, func1.CalledTimes(), 1)
  }
  ```

## Minimal test code

- write trivial tests
  ```go
  func Add(a, b int) int {
    return a + b
  }
  ```
  in test
  ```go
  func Test_SomeOfTwoNumber(...) {
    assert(t, Add(1,2), 3) // It's OK
    assert(t, Add(-1, 2), 1) // trivial bc I don't think there is a high possibility of failure just by calling a standard function.
  }
  ```

### Best Practice

- use test helper

  ```go
  func newItemsForTest() []Item {
    return []Item{
      {Name: "Shampoo", Num: 1},
      {Name: "Pen", Num: 2},
    }
  }

  func Test_1(...) {
    // Arrange
    items := newItemsForTest()
  }

  func Test_2(...) {
    // Arrange
    items := newItemsForTest()
  }
  ```

# Tips

- Alchemy Mock: read `docs/docs/development/Mock.md``
