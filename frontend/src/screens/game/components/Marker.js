import styled from "styled-components";

const MarkerContainer = styled.div`
  height: 80px;
  width: 80px;
  display: flex;
  justify-content: center;
  align-items: center;
`;

const MarkerText = styled.a`
  font-size: 30px;
  font-weight: 800;
  font-family: "Roboto", sans-serif;
  color: ${(props) => props.color};
`;

export default function Marker({ text, color }) {
  return (
    <MarkerContainer>
      <MarkerText color={color}>{text}</MarkerText>
    </MarkerContainer>
  );
}
