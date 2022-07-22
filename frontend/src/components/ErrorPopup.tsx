import { Button } from "@mui/material";
import { useEffect, useState } from "react";
import { SetterOrUpdater } from "recoil";
import styled from "styled-components";

const ErrorContainer = styled.div`
  height: 10rem;
  width: 15rem;
  border-radius: 15px;
  display: flex;
  font-size: 1rem;
  border: 1px solid #dedede;
  flex-direction: column;
  justify-content: space-between;
  align-items: center;
  background-color: white;
  bottom: 20px;
  right: 42rem;
  z-index: 10;
  position: absolute;
  box-shadow: 15px 10px 10px #dedede;
  animation-name: displaceContent;
  animation-duration: 1s;
  animation-delay: 0s;
  animation-iteration-count: 1;
  animation-fill-mode: forwards;
  @keyframes displaceContent {
    from {
      transform: translateX(60rem);
    }
    to {
      transform: translateX(40rem);
    }
  }
`;

export const ErrorPopup = ({
  err,
  setOpen,
}: {
  err: string;
  setOpen: SetterOrUpdater<string | null>;
}) => {


  useEffect(() => {
    setTimeout(() => {
      setOpen(null)
    }, 3000)
  })

  return (
    <ErrorContainer>
      <div style={{ margin: "2rem" }}>{err}</div>
      <Button onClick={() => setOpen(null)}>Close</Button>
    </ErrorContainer>
  );
};
