import { StarsBackground } from "@/components/ui/stars-background";
import React from "react";
import { ShootingStars } from "@/components/ui/shooting-stars";
import { FaEnvelope, FaLock } from "react-icons/fa"; // Icons for inputs

function Login() {
  return (
    <div className="relative bg-black w-full min-h-screen flex items-center justify-center login">
      <StarsBackground className="z-0" />
      <ShootingStars />
      <div className="slide-up relative z-10 bg-black/50 backdrop-blur-md border-2 border-white/50 w-[400px] p-8 rounded-lg text-white shadow-lg transition-transform transform hover:scale-105">
        <h1 className="font-semibold text-3xl mb-1">Hey there,</h1>
        <h1 className="mb-6">Login to get back to your coding!</h1>
        <form className="mt-5">
          <div className="mb-4 relative">
            <label htmlFor="email" className="block font-bold mb-2">
              Email
            </label>
            <div className="relative">
              <FaEnvelope className="absolute left-3 top-3 text-white/70 pointer-events-none" />
              <input
                name="email"
                type="email"
                placeholder="Enter your email.."
                className="w-full pl-10 rounded-md p-2 mb-1 border border-white/80 bg-transparent text-white placeholder-white/70 focus:outline-blue-500 transition"
              />
            </div>
          </div>

          <div className="mb-6 relative">
            <label htmlFor="password" className="block font-bold mb-2">
              Password
            </label>
            <div className="relative">
              <FaLock className="absolute left-3 top-3 text-white/70 pointer-events-none" />
              <input
                name="password"
                type="password"
                placeholder="Enter your password.."
                className="w-full pl-10 rounded-md p-2 mb-1 border border-white/80 bg-transparent text-white placeholder-white/70 focus:outline-blue-500 transition"
              />
            </div>
          </div>

          <button
            type="submit"
            className="cursor-pointer w-full p-2 bg-blue-700 hover:scale-105 hover:bg-white text-black rounded-md hover:bg-blue-800  transition-all duration-300"
          >
            Login
          </button>
        </form>
      </div>
    </div>
  );
}

export default Login;
