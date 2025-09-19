import React, { useEffect, useState, useContext, createContext, useRef, useCallback } from 'react'
import { fetchData } from '../services/backendApi';
export interface User {
  id: number;
  FirstName: string;
  LastName: string;
  Username: string;
  Email: string;
  Password: string;
  CreatedAt: string;
  UpdatedAt: string;
}

export interface Room {
  id: number;
  Name: string;
  RoomID: string;
  CreatedAt: number;
  CreatedBy: string;
  Creator: User;
  Status: string;
}


export interface CreateRoomRequest {
  name: string;
}


export type RoomContextType = {
  createRoom: (roomData: CreateRoomRequest) => Promise<CreateRoomRequest | undefined>;
  getRooms: () => Promise<Room[] | undefined>;
  getRoomById: (id: number) => Promise<Room | undefined>;
};

export const RoomContext = createContext<RoomContextType | undefined>(
  undefined
);


export const RoomProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const getRooms = async (): Promise<Room[] | undefined> => {
    try {
      const response = await fetchData.get("/rooms");
      if (response.status === 200) {
        console.log("Rooms fetched successfully!", response.data.rooms);
        return response.data.rooms as Room[];
      }
    } catch (err: string | any) {
      console.log("Error happened at getRooms function in roomContext file!", err.message);
    }
  }

  const getRoomById = async (id: number): Promise<Room | undefined> => {
    try {
      const response = await fetchData.get(`/rooms/${id}`);
      if (response.status === 200) {
        return response.data.room as Room;
      }
    } catch (err: string | any) {
      console.log("Error happened at getRoomById function in roomContext file!", err.message);

    }
  }

  return (
    <RoomContext.Provider value={{ getRooms, getRoomById }}>
      {children}
    </RoomContext.Provider>
  )
}

export const useRoom = () => {
  return useContext(RoomContext);
}