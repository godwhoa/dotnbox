import styled from "styled-components";

const Dot = styled.div`
  width: 20px;
  height: 20px;
  background-color: ${(props) => props.color || "black"};
`;

export default Dot;
