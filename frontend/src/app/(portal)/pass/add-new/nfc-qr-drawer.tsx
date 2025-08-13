"use client";

import * as React from "react";
import { Button } from "@/components/ui/button";
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerDescription,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from "@/components/ui/drawer";
import { QRCode } from "@/components/ui/qr-code";
import { Usb } from "lucide-react";

interface NfcQrDrawerProps {
  serverNfcData: string;
  serverQrData: string;
}

export function NfcQrDrawer({ serverNfcData, serverQrData }: NfcQrDrawerProps) {
  return (
    <Drawer>
      <DrawerTrigger asChild>
        <Button className="w-full">View Pass</Button>
      </DrawerTrigger>
      <DrawerContent>
        <div className="mx-auto w-full max-w-md">
          <DrawerHeader>
            <DrawerTitle>View NFC & QR Data</DrawerTitle>
            <DrawerDescription>
              These are the latest data received from the server.
            </DrawerDescription>
          </DrawerHeader>

          <div className="flex flex-col sm:flex-row gap-6 p-4">
            {/* NFC Card */}
            <div className="flex-1 rounded-lg border border-gray-300 p-6 shadow-sm flex flex-col">
              <div className="flex items-center space-x-3 mb-4">
                <Usb className="h-6 w-6 text-orange-600" />
                <div>
                  <h3 className="text-lg font-semibold">Write to nfc</h3>
                </div>
              </div>
            </div>

            {/* QR Code Card */}
            <div className="flex-1 rounded-lg border border-gray-300 p-6 shadow-sm flex flex-col items-center">
              <h3 className="mb-4 text-lg font-semibold">QR Code</h3>
              <QRCode
                className="w-48 h-48 rounded border bg-white p-4 shadow"
                data={serverQrData}
              />
              <p className="mt-3 text-center text-sm text-gray-500">
                QR code generated from the server data.
              </p>
            </div>
          </div>

          <DrawerFooter className="flex justify-end space-x-2">
            <DrawerClose asChild>
              <Button variant="outline">Close</Button>
            </DrawerClose>
          </DrawerFooter>
        </div>
      </DrawerContent>
    </Drawer>
  );
}
