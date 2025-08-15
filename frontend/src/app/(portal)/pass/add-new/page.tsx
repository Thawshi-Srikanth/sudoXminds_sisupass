"use client";

import { useState } from "react";
import { ExpDatePicker } from "@/components/portal/pass/exp-date-picker";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { gql, useMutation } from "@apollo/client";
import PassCard, { PassDataState } from "./pass-card";
import client from "@/lib/apolloClient";

const CREATE_PASS_WALLET = gql`
  mutation CreatePassWallet {
    createPassWallet {
      wallet {
        id
      }
    }
  }
`;

const CREATE_PASS = gql`
  mutation CreatePass(
    $categoryId: Int!
    $fromLocationId: Int
    $toLocationId: Int
    $allowedLocationIds: [Int]
    $startDate: Date
    $endDate: Date
  ) {
    createPass(
      categoryId: $categoryId
      fromLocationId: $fromLocationId
      toLocationId: $toLocationId
      allowedLocationIds: $allowedLocationIds
      startDate: $startDate
      endDate: $endDate
    ) {
      passDetails {
        id
      }
    }
  }
`;

export default function Home() {
  const [name, setName] = useState("");
  const [expiry, setExpiry] = useState<Date | undefined>();
  const [selectedPassData, setSelectedPassData] = useState<PassDataState>({});

  const [createWallet, { loading: walletLoading, error: walletError }] =
    useMutation(CREATE_PASS_WALLET, { client });

  const [createPass, { loading: passLoading, error: passError }] = useMutation(
    CREATE_PASS,
    { client }
  );

  const handleCreatePass = async () => {
    try {
      // 1️⃣ Create the pass wallet
      const walletRes = await createWallet();
      const walletId = walletRes.data.createPassWallet.wallet.id;

      // 2️⃣ Loop through each selected pass
      for (const passName in selectedPassData) {
        const pass = selectedPassData[passName];

        await createPass({
          variables: {
            categoryId: pass.categoryId, // Ensure correct type
            fromLocationId: pass.from ? parseInt(pass.from) : null,
            toLocationId: pass.to ? parseInt(pass.to) : null,
            allowedLocationIds: pass.locations?.map((id) => parseInt(id)) || [],
            startDate: new Date().toISOString().split("T")[0], // Example start date
            endDate: expiry ? expiry.toISOString().split("T")[0] : null,
          },
        });
      }

      alert("Passes created successfully!");
      setSelectedPassData({});
      setName("");
      setExpiry(undefined);
    } catch (err) {
      console.error(err);
      alert("Error creating passes");
    }
  };

  return (
    <>
      <div className="col-span-6 flex justify-between w-full px-2">
        <div className="flex justify-between flex-col">
          <h1 className="scroll-m-20 text-2xl font-bold text-balance">
            Add New Pass
          </h1>
          <p className="text-base tracking-wide">
            Fill out the form to create a new pass.
          </p>
        </div>
      </div>

      <div className="col-span-6 flex justify-between w-full px-2">
        <div className="grid grid-cols-4 w-full gap-3">
          <div className="col-span-4 grid w-full items-center gap-3">
            <Label htmlFor="name">Name</Label>
            <Input
              type="text"
              id="name"
              placeholder="Name"
              value={name}
              onChange={(e) => setName(e.target.value)}
            />
          </div>

          <div className="col-span-2 grid w-full items-center gap-3">
            <ExpDatePicker
              value={expiry ? expiry.toISOString().split("T")[0] : undefined}
              onChange={(val: string) =>
                setExpiry(val ? new Date(val) : undefined)
              }
            />
          </div>
        </div>
      </div>

      <div className="col-span-6 flex flex-1 h-full w-full">
        <PassCard
          passData={selectedPassData}
          setPassData={setSelectedPassData}
          onCreatePass={handleCreatePass}
        />
      </div>

      {(walletLoading || passLoading) && <p>Creating pass...</p>}
      {(walletError || passError) && (
        <p>Error creating pass: {walletError?.message || passError?.message}</p>
      )}
    </>
  );
}
