import styled from "styled-components";

const Line = styled.div`
  width: ${(props) => (props.vertical ? "20px" : "80px")};
  height: ${(props) => (props.vertical ? "80px" : "20px")};
  background-color: ${(props) => props.color};
`;

export default Line;
