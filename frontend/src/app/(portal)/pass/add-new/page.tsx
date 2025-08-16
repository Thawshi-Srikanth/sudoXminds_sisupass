"use client";

import { useState } from "react";
import { ExpDatePicker } from "@/components/portal/pass/exp-date-picker";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { gql, useMutation } from "@apollo/client";
import PassCard, { PassDataState } from "./pass-card";
import client from "@/lib/apolloClient";
import { toast } from "sonner";
import { useRouter } from "next/navigation";

const CREATE_PASS_WALLET_WITH_DETAILS = gql`
  mutation CreatePassWalletWithDetails(
    $name: String!
    $expDate: String
    $passDetails: [PassDetailInput!]!
  ) {
    createPassWalletWithDetails(
      name: $name
      expDate: $expDate
      passDetails: $passDetails
    ) {
      wallet {
        walletId
        balance
      }
      createdPassDetails {
        id
      }
    }
  }
`;

export default function Home() {
  const [passName, setPassName] = useState<string>("");
  const [expiry, setExpiry] = useState<Date | undefined>();
  const [selectedPassData, setSelectedPassData] = useState<PassDataState>({});
  const router = useRouter();

  const [createPassWalletWithDetails, { loading, error }] = useMutation(
    CREATE_PASS_WALLET_WITH_DETAILS,
    { client }
  );

  const buildPassDetailsInput = () => {
    return Object.entries(selectedPassData).map(([passName, pass]) => ({
      categoryId: pass.categoryId.toString(),
      fromLocationId: pass.from ? parseInt(pass.from) : null,
      toLocationId: pass.to ? parseInt(pass.to) : null,
      allowedLocationIds: pass.locations?.map((id) => parseInt(id)) || [],
    }));
  };

  const handleCreatePass = async () => {
    try {
      const passDetailsInput = buildPassDetailsInput();

      const { data } = await createPassWalletWithDetails({
        variables: {
          name: passName,
          expDate: expiry ? expiry.toISOString().split("T")[0] : null,
          passDetails: passDetailsInput,
        },
      });

      // Reset state
      setSelectedPassData({});
      setExpiry(undefined);
      setPassName("");

      toast.success("Pass Created Successfully");
      router.push("/pass");
    } catch (err) {
      console.error(err);
      toast.error("Error creating passes");
    }
  };

  return (
    <div className="grid grid-cols-6 gap-x-4 gap-y-6 p-6 w-full h-full">
      <div className="col-span-6 flex justify-between w-full px-2">
        <div className="flex flex-col">
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
              value={passName}
              onChange={(e) => setPassName(e.target.value)}
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
    </div>
  );
}
