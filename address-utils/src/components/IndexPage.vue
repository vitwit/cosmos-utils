<template>
  <div class="container">
    <ul class="nav nav-tabs" id="myTab" role="tablist">
      <li class="nav-item" role="presentation">
        <button
          class="nav-link active"
          id="home-tab"
          data-bs-toggle="tab"
          data-bs-target="#home"
          type="button"
          role="tab"
          aria-controls="home"
          aria-selected="true"
        >
          Valcons | Valconspub
        </button>
      </li>
      <li class="nav-item" role="presentation">
        <button
          class="nav-link"
          id="profile-tab"
          data-bs-toggle="tab"
          data-bs-target="#profile"
          type="button"
          role="tab"
          aria-controls="profile"
          aria-selected="false"
        >
          Account Address
        </button>
      </li>
    </ul>
    <div class="tab-content" id="myTabContent">
      <div
        class="tab-pane fade show active"
        id="home"
        role="tabpanel"
        aria-labelledby="home-tab"
      >
        <div class="card">
          <!-- <div class="card-header">Featured</div> -->
          <div class="card-body">
            <h5 class="card-title text-center">
              Derive valconspub and valcons from base64 encoded pubkey
            </h5>
            <!-- <textarea
              class="form-control"
              placeholder="Please enter pubkey "
              id="floatingTextarea2"
              v-model="validatorPubkey"
              style="height: 100px"
            ></textarea> -->

            <input
              class="form-control form-control-lg"
              disabled
              type="text"
              placeholder="type:tendermint/PubKeyEd25519 "
              aria-label="ValPubkey Type"
            />
            <label
              for="exampleFormControlTextarea1"
              class="form-label text-start"
            ></label>
            <input
              id="exampleFormControlTextarea1"
              class="form-control form-control-lg"
              type="text"
              v-model="valPubkeyValue"
              placeholder="Please enter value"
              aria-label="value"
            />

            <br />
            <p class="text-start">Example:</p>
            <pre class="text-start text-muted">{{
              jsonFormat({
                type: "tendermint/PubKeyEd25519",
                value: "iARPJIXVwen//D6qB5CoQT1KrTK7ffGOkIstFF3KIgk=",
              })
            }}</pre>
            <button class="btn btn-primary text-start" @click="derive()">
              Generate valcons | valconspub
            </button>
            <!-- <button class="btn btn-danger" @click="derive()">
              Clear
            </button> -->
          </div>

          <div v-if="generatValKeys">
            <div class="card">
              <div
                class="card-header"
                style="background-color: green; color: white"
              >
                Generated valcons | valconspub
              </div>
              <div class="card-body">
                <p class="text-start">
                  valconspub : <b>{{ generatedValKeys.valconspub }}</b>
                </p>
                <p class="text-start">
                  valcons : <b>{{ generatedValKeys.valcons }}</b>
                </p>
              </div>
            </div>
          </div>

          <!-- <div class="card-footer text-muted">2 days ago</div> -->
        </div>
      </div>
      <div
        class="tab-pane fade"
        id="profile"
        role="tabpanel"
        aria-labelledby="profile-tab"
      >
        <div class="card">
          <!-- <div class="card-header">Featured</div> -->
          <div class="card-body">
            <h5 class="card-title text-center">
              Deriving a Cosmos Hub address from a public key
            </h5>
            <div class="form">
              <!-- <input
                class="form-control"
                placeholder="Please enter pubkey "
                id="floatingTextarea2"
                style="height: 100px"
                v-model="accountPubkey"
              /> -->
              <input
                class="form-control form-control-lg"
                disabled
                type="text"
                placeholder="type:tendermint/PubKeySecp256k1 "
                aria-label=".form-control-lg example"
              />
              <label
                for="exampleFormControlTextarea1"
                class="form-label text-start"
              ></label>
              <input
                id="exampleFormControlTextarea1"
                class="form-control form-control-lg"
                type="text"
                v-model="accountPubkeyValue"
                placeholder="Please enter value"
                aria-label="value"
              />
            </div>
            <br />
            <p class="text-start">Example:</p>
            <pre class="text-start text-muted">{{
              jsonFormat({
                type: "tendermint/PubKeySecp256k1",
                value: "AtQaCqFnshaZQp6rIkvAPyzThvCvXSDO+9AzbxVErqJP",
              })
            }}</pre>
            <a
              href="#"
              @click="deriveAccountAddress()"
              class="btn btn-primary text-start"
              >Derive Account Address
            </a>
          </div>

          <div v-if="generatAccountAddressKeys">
            <div class="card">
              <div
                class="card-header"
                style="background-color: green; color: white"
              >
                Generated address
              </div>
              <div class="card-body">
                <p class="text-start">
                  address : <b>{{ generatedAccountAddressKey }}</b>
                </p>
              </div>
            </div>
          </div>
          <!-- <div class="card-footer text-muted">2 days ago</div> -->
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { sha256 } from "@cosmjs/crypto";
import { Bech32, fromBase64 } from "@cosmjs/encoding";
import { encodeBech32Pubkey } from "@cosmjs/launchpad";
import { pubkeyToAddress } from "@cosmjs/amino";

export default {
  name: "IndexPage",
  data() {
    return {
      validatorPubkey: {
        type: "tendermint/PubKeyEd25519",
        value: "iARPJIXVwen//D6qB5CoQT1KrTK7ffGOkIstFF3KIgk=",
      },
      valPubkeyValue: "",
      generatValKeys: false,
      generatedValKeys: {
        valconspub:
          "regen:valconspub1zcjduepq3qzy7fy96hq7nllu864q0y9ggy754tfjhd7lrr5s3vk3ghw2ygyspmhej4",
        valcons: "valcons14y3uv3g3fp5k473qtdenmn5cv89y2s5nz7cshu",
      },
      accountPubkey: {
        type: "tendermint/PubKeySecp256k1",
        value: "AtQaCqFnshaZQp6rIkvAPyzThvCvXSDO+9AzbxVErqJP",
      },
      accountPubkeyValue: "",
      generatAccountAddressKeys: false,
      generatedAccountAddressKey: "",
    };
  },
  methods: {
    derive() {
      const pubkey = {
        type: "tendermint/PubKeyEd25519",
        value: this.valPubkeyValue,
      };
      const bech32Pubkey = encodeBech32Pubkey(pubkey, "regen:valconspub");
      console.log("bech32Pubkey ", bech32Pubkey); // regen:valconspub1zcjduepq3qzy7fy96hq7nllu864q0y9ggy754tfjhd7lrr5s3vk3ghw2ygyspmhej4

      const ed25519PubkeyRaw = fromBase64(pubkey.value);
      const addressData = sha256(ed25519PubkeyRaw).slice(0, 20);
      const bech32Address = Bech32.encode("regen:valcons", addressData);
      console.log("bech32Address ", bech32Address); // regen:valcons14y3uv3g3fp5k473qtdenmn5cv89y2s5nz7cshu
      this.generatValKeys = true;
      this.generatedValKeys.valconspub = bech32Pubkey;
      this.generatedValKeys.valcons = bech32Address;
    },
    deriveAccountAddress() {
      const pubkey = {
        type: "tendermint/PubKeySecp256k1",
        value: this.accountPubkeyValue,
      };
      const address = pubkeyToAddress(pubkey, "cosmos");
      console.log({ address });
      this.generatAccountAddressKeys = true;
      this.generatedAccountAddressKey = address;
    },
    jsonFormat(value) {
      return JSON.stringify(value, null, 2);
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.container {
  margin-top: 60px;
}
</style>
