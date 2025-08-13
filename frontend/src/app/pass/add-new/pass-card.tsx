import { useState } from "react";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
} from "@/components/ui/card";

import { ComboBox } from "@/components/ui/combobox";
import { MultiSelect } from "@/components/ui/multi-select";
import { ScrollArea } from "@/components/ui/scroll-area";
import { NfcQrDrawer } from "./nfc-qr-drawer";

// Separate options for each pass type
const busStops = [
  { value: "bus-stop-1", label: "Bus Stop 1" },
  { value: "bus-stop-2", label: "Bus Stop 2" },
];

const trainStations = [
  { value: "station-1", label: "Station 1" },
  { value: "station-2", label: "Station 2" },
];

const libraryLocations = [
  { value: "lib-1", label: "Central Library" },
  { value: "lib-2", label: "Community Library" },
];

const foodLocations = [
  { value: "cafeteria-1", label: "Main Cafeteria" },
  { value: "canteen-1", label: "East Wing Canteen" },
];

const passTypes = ["Bus Pass", "Train Pass", "Library Pass", "Food Pass"];

export default function PassCard() {
  const [selectedPasses, setSelectedPasses] = useState<string[]>([]);
  const [passData, setPassData] = useState<Record<string, any>>({});

  const handleSelectPass = (pass: string) => {
    if (!selectedPasses.includes(pass)) {
      setSelectedPasses((prev) => [...prev, pass]);
    }
  };

  const handleDataChange = (pass: string, field: string, value: any) => {
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
            <div className="flex flex-col h-full divide-y w-full  gap-4">
              {selectedPasses.map((pass) => (
                <div key={pass} className="space-y-3 pb-3">
                  <h4 className="font-medium">{pass}</h4>

                  {pass === "Bus Pass" && (
                    <div className="grid grid-cols-2 gap-3 relative">
                      <ComboBox
                        placeholder="From"
                        options={busStops}
                        onChange={(val) => handleDataChange(pass, "from", val)}
                      />
                      <ComboBox
                        placeholder="To "
                        options={busStops}
                        onChange={(val) => handleDataChange(pass, "to", val)}
                      />
                    </div>
                  )}

                  {pass === "Train Pass" && (
                    <div className="grid grid-cols-2 gap-3 relative">
                      <ComboBox
                        placeholder="From"
                        options={trainStations}
                        onChange={(val) => handleDataChange(pass, "from", val)}
                      />
                      <ComboBox
                        placeholder="To "
                        options={trainStations}
                        onChange={(val) => handleDataChange(pass, "to", val)}
                      />
                    </div>
                  )}

                  {pass === "Library Pass" && (
                    <MultiSelect
                      options={libraryLocations}
                      onValueChange={(vals) =>
                        handleDataChange(pass, "locations", vals)
                      }
                      defaultValue={passData[pass]?.locations || []}
                    />
                  )}

                  {pass === "Food Pass" && (
                    <MultiSelect
                      options={foodLocations}
                      onValueChange={(vals) =>
                        handleDataChange(pass, "locations", vals)
                      }
                      defaultValue={passData[pass]?.locations || []}
                    />
                  )}
                </div>
              ))}
            </div>
          </ScrollArea>
        </CardContent>

        <CardFooter className="px-0">
          <NfcQrDrawer
            serverNfcData="https://sisupass.com/pass/123s"
            serverQrData="https://sisupass.com/pass/123"
          />
        </CardFooter>
      </Card>
    </div>
  );
}
