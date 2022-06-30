import styled from "styled-components";

const Button = styled.button`
  background-color: #8a8a8a;
  border: none;
  color: white;
  padding: 1rem 1rem;
  text-align: center;
  display: flex;
  align-items: center;
  justify-content: start;
  font-size: 1rem;
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
  font-size: 1rem;
  border: 1px solid #dedede;
  flex-direction: column;
  justify-content: space-between;
  align-items: center;
  background-color: white;
  top: 20px;
  z-index:10;
  position: absolute;
  box-shadow: 15px 10px 10px #dedede;
  animation-name            : displaceContent;
  animation-duration        : 1s;
  animation-delay           : 0s;
  animation-iteration-count : 1;
  animation-fill-mode       : forwards;
  @keyframes displaceContent {
    from { transform : translateX(20rem) }
    to   { transform : translateX(0rem) } 
  }
`;

export const ErrorPopup = ({ err, setOpen }: { err: string, setOpen(arg: boolean):void }) => {
  return (
        <ErrorContainer>
            <div style={{margin:'2rem'}}>{err}</div>
          <Button onClick={() => setOpen(false)}>Close</Button>
        </ErrorContainer>
  );
};