import { Button } from "@/components/ui/button";
import { useGoogleLogin } from "@react-oauth/google";

const GoogleSignInButton = ({ handleGoogleSignIn }) => {
  const login = useGoogleLogin({
    onSuccess: (tokenResponse) => handleGoogleSignIn(tokenResponse),
  });
  return (
    <Button size="xl" variant="outline" onClick={() => login()}>
      Sign In Google
    </Button>
  );
};
export default GoogleSignInButton;
