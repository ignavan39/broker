import { useState } from "react";
import styled from "styled-components";

const Button = styled.button`
  background-color: #8a8a8a;
  border: none;
  color: white;
  padding: 1rem 2rem;
  text-align: center;
  display: flex;
  align-items: center;
  justify-content: start;
  font-size: 18px;
  min-width: 2rem;
  margin: 0.2rem 0;
  border-radius: 10px;
  &:hover {
    background-color: #ff2626;
  }
`;

const ErrorContainer = styled.div`
  height: 10rem;
  width: 20rem;
  border-radius: 15px;
  display: flex;
  font-size: 2rem; 2rem;
  flex-direction: column;
  justify-content: space-between;
  align-items: center;
  background-color: #dedede;
  top: 0;
  position: absolute;
`;

export const ErrorPopup = ({ err }: { err: string }) => {
  const [open, close] = useState(true);
  return (
    <>
      {open ? (
        <ErrorContainer>
            <div style={{margin:'2rem'}}>{err}</div>
          <Button onClick={() => close(false)}>Close</Button>
        </ErrorContainer>
      ) : (
        <></>
      )}
    </>
  );
};
