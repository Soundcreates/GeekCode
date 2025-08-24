import React from "react";
import { FloatingDock } from "../components/ui/floating-docs";
import { StarsBackground } from "../components/ui/stars-background";
import { TextGenerateEffect } from "../components/ui/text-generate-effect";

import { useNavigate } from "react-router";
import { HomeIcon, Apple, ContactIcon } from "lucide-react";

function Landing() {
  const navigate = useNavigate();

  return (
    <div className="relative w-full h-screen bg-black text-white overflow-x-auto">
      {/* Stars background */}
      <StarsBackground className="z-0" />

      {/* Rest of the page */}
      <div className="relative z-10 h-screen  flex flex-col items-center justify-center gap-5 p-2">
        <div>
          <FloatingDock
            items={[
              {
                title: "Home",
                icon: <HomeIcon />,
                href: "/Login",
              },
              {
                title: "About",
                icon: <Apple />,
                href: "/get-started",
              },
              {
                title: "Contact",
                icon: <ContactIcon />,
                href: "/contact",
              },
            ]}
          />
        </div>

        <div className="flex flex-col items-center justify-center  h-full">
          <TextGenerateEffect words="Welcome to GeekCode" />
          <p className="mt-6 text-lg text-white/70 max-w-xl text-center generate-p">
            Dive into the world of collaborative coding with GeekCode. Create or
            join coding rooms and code together in real-time!
          </p>

          <div className="mt-5 flex gap-10 justify-between">
            <div className="p-2 bg-black rounded-md text-3xl transition-transform duration-300 hover:scale-110 cursor-pointer">
              <h1
                className="text-white/70 transition-all duration-300 hover:text-white hover:drop-shadow-[0_0_12px_rgba(255,255,255,1)]"
                onClick={() => navigate("/login")}
              >
                Login
              </h1>
            </div>

            <div className="p-2 bg-black rounded-md  text-3xl transition-transform duration-300 hover:scale-110 cursor-pointer">
              <h1
                className="text-white/70 transition-all duration-300 hover:text-white hover:drop-shadow-[0_0_12px_rgba(255,255,255,1)]"
                onClick={() => navigate("/login")}
              >
                Get Started
              </h1>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Landing;
