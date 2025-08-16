import { TransactionList } from "./transaction-list";

export default function Home() {
  return (
    <div className="grid grid-cols-6 gap-x-4 gap-y-6 p-6 w-full h-full">
      <div className="col-span-6 flex justify-between flex-col px-2">
        <h1 className="scroll-m-20  text-2xl font-bold text-balance">
          Transaction
        </h1>
        <p className=" text-base tracking-wide">Manage your transactions with ease.</p>
      </div>

      <div className="col-span-6 flex flex-1 justify-between flex-col">
        <TransactionList />
      </div>
    </div>
  );
}
