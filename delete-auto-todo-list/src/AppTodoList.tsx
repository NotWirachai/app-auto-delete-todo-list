import React, { useState } from "react";
import { Todo } from "./AppTodoList.interface";
import { Container, Column, Button, ColumnTable } from "./AppTodoList.style";

const initialTodos: Todo[] = [
  { type: "Fruit", name: "Apple" },
  { type: "Vegetable", name: "Broccoli" },
  { type: "Vegetable", name: "Mushroom" },
  { type: "Fruit", name: "Banana" },
  { type: "Vegetable", name: "Tomato" },
  { type: "Fruit", name: "Orange" },
  { type: "Fruit", name: "Mango" },
  { type: "Fruit", name: "Pineapple" },
  { type: "Vegetable", name: "Cucumber" },
  { type: "Fruit", name: "Watermelon" },
  { type: "Vegetable", name: "Carrot" },
];

const TodoApp: React.FC = () => {
  const [mainList, setMainList] = useState<Todo[]>(initialTodos);
  const [movedList, setMovedList] = useState<{ [key: string]: Todo[] }>({
    Fruit: [],
    Vegetable: [],
  });

  const moveItem = (item: Todo) => {
    setMainList((prev) => prev.filter((i) => i.name !== item.name));
    setMovedList((prev) => ({
      ...prev,
      [item.type]: [...prev[item.type], item],
    }));

    setTimeout(() => {
      setMovedList((prev) => ({
        ...prev,
        [item.type]: prev[item.type].filter((i) => i.name !== item.name),
      }));
      setMainList((prev) => [...prev, item]);
    }, 5000);
  };

  const moveBack = (item: Todo) => {
    setMovedList((prev) => ({
      ...prev,
      [item.type]: prev[item.type].filter((i) => i.name !== item.name),
    }));
    setMainList((prev) => [...prev, item]);
  };

  return (
    <Container>
      <Column>
        {mainList.map((item) => (
          <Button key={item.name} onClick={() => moveItem(item)}>
            {item.name}
          </Button>
        ))}
      </Column>
      <ColumnTable>
        <div className="column-title">
          <h3>Fruits</h3>
        </div>
        {movedList.Fruit.map((item) => (
          <Button key={item.name} onClick={() => moveBack(item)} stay>
            {item.name}
          </Button>
        ))}
      </ColumnTable>
      <ColumnTable>
        <div className="column-title">
          <h3>Vegetable</h3>
        </div>
        {movedList.Vegetable.map((item) => (
          <Button key={item.name} onClick={() => moveBack(item)} stay>
            {item.name}
          </Button>
        ))}
      </ColumnTable>
    </Container>
  );
};

export default TodoApp;
