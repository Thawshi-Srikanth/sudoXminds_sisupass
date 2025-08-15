"use client";

import { useState } from "react";
import { useQuery, gql } from "@apollo/client";
import client from "@/lib/apolloClient";

import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

import { ComboBox } from "@/components/ui/combobox";
import { MultiSelect } from "@/components/ui/multi-select";
import { ScrollArea } from "@/components/ui/scroll-area";
import { NfcQrDrawer } from "./nfc-qr-drawer";

export const GET_PASS_DATA = gql`
  query GetPassData {
    allLocationTypes {
      id
      name
    }
    allLocations {
      id
      name
      locationType {
        id
        name
      }
    }
    allPassCategories {
      id
      name
      allowedLocationTypes {
        id
        name
      }
    }
  }
`;

export interface Location {
  id: string;
  name: string;
  locationType: {
    id: string;
    name: string;
  };
}

export interface LocationType {
  id: string;
  name: string;
}

export interface PassCategory {
  id: string;
  name: string;
  allowedLocationTypes: LocationType[];
}

export interface PassDataQuery {
  allLocationTypes: LocationType[];
  allLocations: Location[];
  allPassCategories: PassCategory[];
}

export interface PassDataState {
  [passName: string]: {
    categoryId: number;
    from?: string;
    to?: string;
    locations?: string[];
  };
}

interface PassCardProps {
  passData: PassDataState;
  setPassData: React.Dispatch<React.SetStateAction<PassDataState>>;
  onCreatePass: () => void;
}

export default function PassCard({
  passData,
  setPassData,
  onCreatePass,
}: PassCardProps) {
  const { data, loading, error } = useQuery<PassDataQuery>(GET_PASS_DATA, {
    client,
  });
  const [selectedPasses, setSelectedPasses] = useState<string[]>([]);

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error loading data</div>;

  // Map locationType name to its locations
  const locationOptions: Record<string, { value: string; label: string }[]> =
    {};
  data?.allLocationTypes.forEach((lt) => {
    locationOptions[lt.name] =
      data.allLocations
        .filter((loc) => loc.locationType.name === lt.name)
        .map((loc) => ({ value: loc.id, label: loc.name })) || [];
  });

  const passTypes = data?.allPassCategories.map((cat) => cat.name) || [];

  const handleSelectPass = (pass: string) => {
    if (!selectedPasses.includes(pass)) {
      setSelectedPasses((prev) => [...prev, pass]);
    }
  };

  const handleDataChange = (
    pass: string,
    field: keyof PassDataState[string],
    value: any
  ) => {
    setPassData((prev) => ({
      ...prev,
      [pass]: {
        ...prev[pass],
        [field]: value,
      },
    }));
  };

  return (
    <div className="col-span-6 flex flex-1 w-full px-2">
      <Card className="w-full h-full shadow-none p-0 border-0">
        <CardHeader className="px-0">
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" className="shadow-none">
                Select Pass
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent className="min-w-sm">
              <DropdownMenuLabel>Pass</DropdownMenuLabel>
              <DropdownMenuSeparator />
              {passTypes.map((pass) => (
                <DropdownMenuItem
                  key={pass}
                  onClick={() => handleSelectPass(pass)}
                >
                  {pass}
                </DropdownMenuItem>
              ))}
            </DropdownMenuContent>
          </DropdownMenu>
        </CardHeader>

        <CardContent className="space-y-4 px-0">
          <ScrollArea className="h-[calc(100vh-500px)] w-full">
            <div className="flex flex-col h-full divide-y w-full gap-4">
              {selectedPasses.map((pass) => {
                const allowedTypes =
                  data?.allPassCategories.find((c) => c.name === pass)
                    ?.allowedLocationTypes || [];

                return (
                  <div key={pass} className="space-y-3 pb-3">
                    <h4 className="font-medium">{pass}</h4>

                    {allowedTypes.map((type) => {
                      const options = locationOptions[type.name] || [];
                      if (pass === "Bus Pass" || pass === "Train Pass") {
                        return (
                          <div
                            key={type.id}
                            className="grid grid-cols-2 gap-3 relative"
                          >
                            <ComboBox
                              placeholder="From"
                              options={options}
                              onChange={(val) =>
                                handleDataChange(pass, "from", val)
                              }
                            />
                            <ComboBox
                              placeholder="To"
                              options={options}
                              onChange={(val) =>
                                handleDataChange(pass, "to", val)
                              }
                            />
                          </div>
                        );
                      } else {
                        return (
                          <MultiSelect
                            key={type.id}
                            options={options}
                            onValueChange={(vals) =>
                              handleDataChange(pass, "locations", vals)
                            }
                            defaultValue={passData[pass]?.locations || []}
                          />
                        );
                      }
                    })}
                  </div>
                );
              })}
            </div>
          </ScrollArea>
        </CardContent>

        <CardFooter className="px-0">
          <Button onClick={onCreatePass} variant="default">
            Create Pass
          </Button>
          {/* <NfcQrDrawer
            serverNfcData="https://sisupass.com/pass/123s"
            serverQrData="https://sisupass.com/pass/123"
          /> */}
        </CardFooter>
      </Card>
    </div>
  );
}
