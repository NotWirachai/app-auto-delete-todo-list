import styled from "styled-components";

export const Container = styled.div`
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 20px;
  padding: 20px;
  height: calc(100vh - 300px);
`;

export const Column = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 10px;
  gap: 10px;
`;

export const ColumnTable = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  border: 1px solid gray;
  gap: 10px;

  .column-title {
    background: gray;
    text-align: center;
    width: 100%;
    color: white;
  }
`;

export const Button = styled.button<{ stay?: boolean }>`
  padding: 10px;
  background: transparent;
  color: white;
  border: 1px solid gray;
  width: ${(props) => (props.stay ? "90%" : "100%")};
  cursor: pointer;
  transition: background 0.3s;
  color: black;

  &:hover {
    background: #f4f4f4;
  }
`;
