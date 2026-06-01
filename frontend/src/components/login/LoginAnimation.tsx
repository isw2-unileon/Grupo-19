import Lottie from "lottie-react";
import animationData from "../../assets/CarroDeCompras.json";

export default function LoginAnimation() {
  return (
    <div style={styles.wrapper}>
      <Lottie animationData={animationData} loop={true} />
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  wrapper: {
    width: "100%",
    height: "100%",
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
  },
};